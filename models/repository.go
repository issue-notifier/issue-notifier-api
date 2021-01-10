package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/issue-notifier/issue-notifier-api/database"
)

// Repository struct to store repository information from database
type Repository struct {
	RepoID      uuid.UUID `json:"repoID" db:"repo_id"`
	RepoName    string    `json:"repoName" db:"repo_name"`
	LastEventAt time.Time `json:"lastEventAt" db:"last_event_at"`
}

// Label struct to store label information from database
type Label struct {
	Name  string `json:"name" db:"label_name"`
	Color string `json:"color" db:"label_color"`
}

// Labels is of type `[]Label`
type Labels []Label

// Value for the `Labels` struct to implement the driver Valuer interface. This method
// enables types to convert themselves to a driver Value.
func (a Labels) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan for the `Labels` struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *Labels) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

// GetAllRepositories gets all repositories from the database
func GetAllRepositories() ([]Repository, error) {
	sqlQuery := `SELECT * FROM GLOBAL_REPOSITORY`

	rows, err := database.DB.Query(sqlQuery)

	var data []Repository

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var repoID uuid.UUID
		var repoName string
		var lastEventAt time.Time
		if err := rows.Scan(&repoID, &repoName, &lastEventAt); err != nil {
			return nil, fmt.Errorf("[GetAllRepositories]: %v", err)
		}

		data = append(data, Repository{
			RepoID:      repoID,
			RepoName:    repoName,
			LastEventAt: lastEventAt,
		})
	}

	return data, nil
}

// GetRepositoryIDByName gets `repoID` for the given `repoName`
func GetRepositoryIDByName(repoName string) (string, error) {
	sqlQuery := `SELECT REPO_ID FROM GLOBAL_REPOSITORY WHERE REPO_NAME=$1`

	var repoID uuid.UUID
	row := database.DB.QueryRow(sqlQuery, repoName)
	err := row.Scan(&repoID)

	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("[GetRepositoryIDByName]: %v", err)
	}

	return repoID.String(), err
}

// CreateRepository creates a new repository in the database
func CreateRepository(repoName string) (string, error) {
	sqlQuery := `INSERT INTO GLOBAL_REPOSITORY (REPO_NAME) VALUES ($1) RETURNING REPO_ID`

	var repoID uuid.UUID
	err := database.DB.QueryRow(sqlQuery, repoName).Scan(&repoID)
	if err != nil {
		return "", fmt.Errorf("[CreateRepository]: %v", err)
	}

	return repoID.String(), nil
}

// UpdateLastEventAtByRepoID updates `lastEventAt` time for the given `repoID`
func UpdateLastEventAtByRepoID(repoID, lastEventAt string) error {
	sqlQuery := `UPDATE GLOBAL_REPOSITORY SET LAST_EVENT_AT = $1 WHERE REPO_ID = $2`

	_, err := database.DB.Exec(sqlQuery, lastEventAt, repoID)
	if err != nil {
		return fmt.Errorf("[UpdateLastEventAtByRepoID]: %v", err)
	}

	return nil
}

// DeleteRepositoriesWithNoLabels deletes repositories which have no subscribed `labels`
func DeleteRepositoriesWithNoLabels() error {
	sqlQuery := "DELETE FROM GLOBAL_REPOSITORY WHERE REPO_ID NOT IN (SELECT USER_SUBSCRIPTION.REPO_ID FROM USER_SUBSCRIPTION)"

	_, err := database.DB.Exec(sqlQuery)
	if err != nil {
		return fmt.Errorf("[DeleteRepositoriesWithNoLabels]: %v", err)
	}

	return nil
}
