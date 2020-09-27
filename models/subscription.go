package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hemakshis/issue-notifier-api/database"
)

type SubscriptionModel struct {
	RepoID               uuid.UUID `json:"repoID" db:"repo_id"`
	UserID               uuid.UUID `json:"userID" db:"user_id"`
	Label                string    `json:"label" db:"label"`
	LastNotificationSent time.Time `json:"lastNotificationSent" db:"last_notification_sent"`
}

func GetSubscriptionsByUserID(userID string) ([]map[string]string, error) {
	sqlQuery := `SELECT DISTINCT GLOBAL_REPOSITORY.REPO_NAME, GLOBAL_REPOSITORY.HTML_URL, USER_SUBSCRIPTION.LABEL 
		FROM GLOBAL_REPOSITORY 
		INNER JOIN USER_SUBSCRIPTION ON GLOBAL_REPOSITORY.REPO_ID = USER_SUBSCRIPTION.REPO_ID 
		AND USER_SUBSCRIPTION.USER_ID=$1`

	rows, err := database.DB.Query(sqlQuery, userID)

	var data []map[string]string

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d map[string]string
		var repoName, htmlURL, label string
		if err := rows.Scan(&repoName, &htmlURL, &label); err != nil {
			return nil, err
		}

		d = map[string]string{
			"repoName": repoName,
			"htmlURL":  htmlURL,
			"label":    label,
		}

		data = append(data, d)
	}

	return data, nil
}

func GetSubscribedLabelsByUserIDAndRepoID(userID, repoID string) ([]string, error) {
	sqlQuery := `SELECT LABEL FROM USER_SUBSCRIPTION WHERE USER_ID=$1 AND REPO_ID=$2`

	rows, err := database.DB.Query(sqlQuery, userID, repoID)

	var labels []string

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var label string
		if err := rows.Scan(&label); err != nil {
			return nil, err
		}

		labels = append(labels, label)
	}

	return labels, nil
}

func CreateSubscription(userID, repoID string, labels []string) error {
	sqlQuery := `INSERT INTO USER_SUBSCRIPTION (USER_ID, REPO_ID, LABEL) VALUES `

	valuesPlaceholder := make([]string, 0)
	values := make([]interface{}, 0)
	for i, l := range labels {
		valuesPlaceholder = append(valuesPlaceholder, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		values = append(values, userID)
		values = append(values, repoID)
		values = append(values, l)
	}

	sqlQuery = sqlQuery + strings.Join(valuesPlaceholder, ",")
	_, err := database.DB.Exec(sqlQuery, values...)

	return err
}

func RemoveSubscription(userID, repoID string, labels []string) error {
	sqlQuery := `DELETE FROM USER_SUBSCRIPTION WHERE USER_ID=$1 AND REPO_ID=$2 AND `

	valuesPlaceholder := make([]string, 0)
	values := make([]interface{}, 0)
	values = append(values, userID)
	values = append(values, repoID)

	for i, l := range labels {
		valuesPlaceholder = append(valuesPlaceholder, fmt.Sprintf("LABEL=$%d", i+3))
		values = append(values, l)
	}

	sqlQuery = sqlQuery + strings.Join(valuesPlaceholder, " OR ")
	_, err := database.DB.Exec(sqlQuery, values...)

	return err
}
