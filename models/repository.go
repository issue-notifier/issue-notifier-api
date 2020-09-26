package models

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hemakshis/issue-notifier-api/database"
)

type Repository struct {
	RepoName string `json:"repoName" db:"repo_name"`
	ApiURL   string `json:"apiUrl" db:"api_url"`
	HtmlURL  string `json:"htmlUrl" db:"html_url"`
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

func CreateRepository(repoName, apiURL, htmlURL string) (string, error) {
	sqlQuery := `INSERT INTO GLOBAL_REPOSITORY (REPO_NAME, API_URL, HTML_URL) VALUES ($1, $2, $3) RETURNING REPO_ID`

	var repoID uuid.UUID
	err := database.DB.QueryRow(sqlQuery, repoName, apiURL, htmlURL).Scan(&repoID)
	if err != nil {
		panic(err)
	}

	return repoID.String(), err
}
