package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/issue-notifier/issue-notifier-api/database"
)

type Repository struct {
	RepoName string `json:"repoName" db:"repo_name"`
}

type Label struct {
	Name  string `json:"name" db:"label_name"`
	Color string `json:"color" db:"label_color"`
}

type Labels []Label

func (a Labels) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Label struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *Labels) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func GetAllRepositories() ([]Repository, error) {
	sqlQuery := `SELECT REPO_NAME FROM GLOBAL_REPOSITORY`

	rows, err := database.DB.Query(sqlQuery)

	var data []Repository

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d Repository
		var repoName string
		if err := rows.Scan(&repoName); err != nil {
			return nil, err
		}

		d.RepoName = repoName

		data = append(data, d)
	}

	return data, nil

}

func GetRepositoryIDByName(repoName string) (string, error) {
	sqlQuery := `SELECT REPO_ID FROM GLOBAL_REPOSITORY WHERE REPO_NAME=$1`

	var repoID uuid.UUID
	row := database.DB.QueryRow(sqlQuery, repoName)
	err := row.Scan(&repoID)

	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	return repoID.String(), err
}

func CreateRepository(repoName string) (string, error) {
	sqlQuery := `INSERT INTO GLOBAL_REPOSITORY (REPO_NAME) VALUES ($1) RETURNING REPO_ID`

	var repoID uuid.UUID
	err := database.DB.QueryRow(sqlQuery, repoName).Scan(&repoID)
	if err != nil {
		panic(err)
	}

	return repoID.String(), err
}

func DeleteRepositoriesWithNoLabels(repoName string) error {
	sqlQuery := "DELETE FROM GLOBAL_REPOSITORY WHERE REPO_ID NOT IN (SELECT USER_SUBSCRIPTION.REPO_ID FROM USER_SUBSCRIPTION)"

	_, err := database.DB.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return err
}
