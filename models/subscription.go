package models

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/issue-notifier/issue-notifier-api/database"
)

// Subscription struct to store subscription info for a single repository
type Subscription struct {
	RepoName string `json:"repoName" db:"repo_name"`
	Labels   Labels `json:"labels" db:"labels"`
}

// UserIDLabel struct stores userID-label pair
type UserIDLabel struct {
	UserID uuid.UUID `json:"userID"`
	Label  string    `json:"label"`
}

// GetSubscriptionsByUserID gets all subscribed `labels` and the `userID` of the user who has subscribed for that particular for the give `repoID`
func GetSubscriptionsByUserID(userID string) ([]Subscription, error) {
	sqlQuery := `SELECT DISTINCT GR.REPO_NAME, US.LABELS 
		FROM GLOBAL_REPOSITORY GR
		INNER JOIN USER_SUBSCRIPTION US 
		ON GR.REPO_ID = US.REPO_ID 
		AND US.USER_ID=$1`

	rows, err := database.DB.Query(sqlQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("[GetSubscriptionsByUserID]: %v", err)
	}
	defer rows.Close()

	var data []Subscription
	for rows.Next() {
		var repoName string
		var labels Labels
		if err := rows.Scan(&repoName, &labels); err != nil {
			return nil, fmt.Errorf("[GetSubscriptionsByUserID]: %v", err)
		}

		data = append(data, Subscription{
			RepoName: repoName,
			Labels:   labels,
		})
	}

	return data, nil
}

// GetSubscriptionsByRepoID gets all subscriptions for the given `repoID`
func GetSubscriptionsByRepoID(repoID string) ([]UserIDLabel, error) {
	sqlQuery := `SELECT USER_ID, LABELS.NAME
		FROM USER_SUBSCRIPTION, JSONB_TO_RECORDSET(USER_SUBSCRIPTION.LABELS) 
		AS LABELS(NAME TEXT, COLOR TEXT) 
		WHERE REPO_ID=$1`

	rows, err := database.DB.Query(sqlQuery, repoID)
	if err != nil {
		return nil, fmt.Errorf("[GetSubscriptionsByRepoID]: %v", err)
	}
	defer rows.Close()

	var data []UserIDLabel
	for rows.Next() {
		var label string
		var userID uuid.UUID
		if err := rows.Scan(&userID, &label); err != nil {
			return nil, fmt.Errorf("[GetSubscriptionsByRepoID]: %v", err)
		}

		data = append(data, UserIDLabel{
			UserID: userID,
			Label:  label,
		})
	}

	return data, nil
}

// GetSubscribedLabelsByUserIDAndRepoID gets all the subscribed `labels` for the given `userID` and `repoID`
func GetSubscribedLabelsByUserIDAndRepoID(userID, repoID string) (Labels, error) {
	sqlQuery := `SELECT LABELS FROM USER_SUBSCRIPTION WHERE USER_ID=$1 AND REPO_ID=$2`
	// TODO sqlQuery := `SELECT US.LABELS FROM USER_SUBSCRIPTION US INNER JOIN GLOBAL_REPOSITORY GR ON US.REPO_ID = GR.REPO_ID WHERE USER_ID=$1 AND REPO_ID=$2`

	row := database.DB.QueryRow(sqlQuery, userID, repoID)

	var labels Labels
	err := row.Scan(&labels)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("[GetSubscribedLabelsByUserIDAndRepoID]: %v", err)
	}

	return labels, err
}

// CreateSubscriptions creates a new entry in the database with given `userID`, `repoID` and `labels`
func CreateSubscriptions(userID, repoID string, labels Labels) error {
	sqlQuery := `INSERT INTO USER_SUBSCRIPTION (USER_ID, REPO_ID, LABELS) VALUES ($1, $2, $3)`

	_, err := database.DB.Exec(sqlQuery, userID, repoID, labels)
	if err != nil {
		return fmt.Errorf("[CreateSubscriptions]: %v", err)
	}

	return nil
}

// UpdateSubscriptions updates an existing subscription by adding the new `labels` for the given `userID` and `repoID`
func UpdateSubscriptions(userID, repoID string, labels Labels) error {
	sqlQuery := `UPDATE USER_SUBSCRIPTION SET LABELS = LABELS || $1 WHERE USER_ID=$2 AND REPO_ID=$3`

	_, err := database.DB.Exec(sqlQuery, labels, userID, repoID)
	if err != nil {
		return fmt.Errorf("[UpdateSubscriptions]: %v", err)
	}

	return nil
}

// RemoveSubscriptions removes all the given `labels` for the given `userID` and `repoID`
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
			return fmt.Errorf("[RemoveSubscriptions]: %v", err)
		}
	}

	// Clean up subscriptions with no labels
	sqlQuery := "DELETE FROM USER_SUBSCRIPTION WHERE JSONB_ARRAY_LENGTH(LABELS) = 0"

	_, err := database.DB.Exec(sqlQuery)
	if err != nil {
		return fmt.Errorf("[RemoveSubscriptions]: %v", err)
	}

	return nil
}
