package database

import (
	"database/sql"
	"fmt"

	"github.com/issue-notifier/issue-notifier-api/utils"

	// Postgres driver for sql
	_ "github.com/lib/pq"
)

// DB to initialize postgres database
var DB *sql.DB

// Init initializes postgres database
func Init(environment string, dbConfigs ...string) {
	var err error
	if environment == "production" {
		DB, err = sql.Open("postgres", dbConfigs[0])
	} else {
		connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbConfigs[0], dbConfigs[1], dbConfigs[2], dbConfigs[3])
		DB, err = sql.Open("postgres", connectionString)
	}

	if err != nil {
		utils.LogError.Fatalln("Failed to connect to the database. Error:", err)
	}
	utils.LogInfo.Println("Successfully connected to the database")
}
