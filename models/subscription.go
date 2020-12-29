package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/issue-notifier/issue-notifier-api/database"
)

type SubscriptionModel struct {
	RepoID               uuid.UUID `json:"repoID" db:"repo_id"`
	UserID               uuid.UUID `json:"userID" db:"user_id"`
	Label                string    `json:"label" db:"label"`
	LastNotificationSent time.Time `json:"lastNotificationSent" db:"last_notification_sent"`
}

func GetSubscriptionsByUserID(userID string) ([]map[string]interface{}, error) {
	sqlQuery := `SELECT DISTINCT GLOBAL_REPOSITORY.REPO_NAME, USER_SUBSCRIPTION.LABELS 
		FROM GLOBAL_REPOSITORY 
		INNER JOIN USER_SUBSCRIPTION ON GLOBAL_REPOSITORY.REPO_ID = USER_SUBSCRIPTION.REPO_ID 
		AND USER_SUBSCRIPTION.USER_ID=$1`

	rows, err := database.DB.Query(sqlQuery, userID)

	var data []map[string]interface{}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var repoName string
		var labels Labels
		if err := rows.Scan(&repoName, &labels); err != nil {
			return nil, err
		}

		data = append(data, map[string]interface{}{
			"repoName": repoName,
			"labels":   labels,
		})
	}

	return data, nil
}

func GetSubscriptionsByRepoID(repoID string) ([]map[string]interface{}, error) {
	sqlQuery := `SELECT USER_ID, LABELS.NAME
		FROM USER_SUBSCRIPTION, JSONB_TO_RECORDSET(USER_SUBSCRIPTION.LABELS) 
		AS LABELS(NAME TEXT, COLOR TEXT) 
		WHERE REPO_ID=$1`

	rows, err := database.DB.Query(sqlQuery, repoID)

	var data []map[string]interface{}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var label string
		var userID uuid.UUID
		if err := rows.Scan(&userID, &label); err != nil {
			return nil, err
		}

		data = append(data, map[string]interface{}{
			"userId": userID,
			"label":  label,
		})
	}

	return data, nil
}

func GetSubscribedLabelsByUserIDAndRepoID(userID, repoID string) (Labels, error) {
	sqlQuery := `SELECT LABELS FROM USER_SUBSCRIPTION WHERE USER_ID=$1 AND REPO_ID=$2`

	row := database.DB.QueryRow(sqlQuery, userID, repoID)

	var labels Labels
	err := row.Scan(&labels)

	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	return labels, err
}

func CreateSubscriptions(userID, repoID string, labels Labels) error {
	sqlQuery := `INSERT INTO USER_SUBSCRIPTION (USER_ID, REPO_ID, LABELS) VALUES ($1, $2, $3)`

	_, err := database.DB.Exec(sqlQuery, userID, repoID, labels)

	return err
}

func UpdateSubscriptions(userID, repoID string, labels Labels) error {
	sqlQuery := `UPDATE USER_SUBSCRIPTION SET LABELS = LABELS || $1 WHERE USER_ID=$2 AND REPO_ID=$3`

	_, err := database.DB.Exec(sqlQuery, labels, userID, repoID)

	return err
}

func RemoveSubscriptions(userID, repoID string, labels Labels) error {
	// TODO: Make it more efficient
	for _, l := range labels {

		sqlQuery := `UPDATE USER_SUBSCRIPTION 
			SET LABELS = LABELS - (
				SELECT I FROM GENERATE_SERIES(0, JSONB_ARRAY_LENGTH(LABELS)-1) AS I 
				WHERE LABELS->I->>'name' = $1
			)::INTEGER
			WHERE USER_ID = $2 AND REPO_ID = $3`

		_, err := database.DB.Exec(sqlQuery, l.Name, userID, repoID)

		if err != nil {
			return err
		}
	}

	// Clean up subscriptions with no labels
	sqlQuery := "DELETE FROM USER_SUBSCRIPTION WHERE JSONB_ARRAY_LENGTH(LABELS) = 0"

	_, err := database.DB.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return nil
}
