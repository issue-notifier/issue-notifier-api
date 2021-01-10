package models

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/issue-notifier/issue-notifier-api/database"
)

// User struct to store user information from database
type User struct {
	UserID   uuid.UUID `json:"userID" db:"user_id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
}

// GetUserIDByUsername gets `userID` for the given `username`
func GetUserIDByUsername(username string) (string, error) {
	sqlQuery := `SELECT USER_ID FROM GITHUB_USER WHERE USERNAME=$1`

	var userID uuid.UUID
	row := database.DB.QueryRow(sqlQuery, username)
	err := row.Scan(&userID)

	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("[GetUserIDByUsername]: %v", err)
	}

	return userID.String(), nil
}

// CreateUser creates a new `User` in the database and returns the `userID` of the newly created `User`
func CreateUser(username, email string) (string, error) {
	sqlQuery := `INSERT INTO GITHUB_USER (USERNAME, EMAIL) VALUES ($1, $2) RETURNING USER_ID`

	var userID uuid.UUID
	err := database.DB.QueryRow(sqlQuery, username, email).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("[CreateUser]: %v", err)
	}

	return userID.String(), nil
}
