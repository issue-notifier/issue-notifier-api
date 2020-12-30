package main

import (
	"log"
	"net/http"
	"os"

	"github.com/issue-notifier/issue-notifier-api/database"
	"github.com/issue-notifier/issue-notifier-api/routes"
	"github.com/issue-notifier/issue-notifier-api/session"
	"github.com/joho/godotenv"
)

// Env vars
var (
	PORT string

	DB_USER string
	DB_PASS string
	DB_NAME string

	SESSION_AUTH_KEY string

	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT = os.Getenv("PORT")
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_NAME = os.Getenv("DB_NAME")
	SESSION_AUTH_KEY = os.Getenv("SESSION_AUTH_KEY")
	GITHUB_CLIENT_ID = os.Getenv("GITHUB_CLIENT_ID")
	GITHUB_CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")

	database.Init(DB_USER, DB_PASS, DB_NAME)
	defer database.DB.Close()

	routes.Init(GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET)

	session.Init(SESSION_AUTH_KEY)

	log.Fatal(http.ListenAndServe(PORT, routes.Router))
}

// TODOs MAJOR: Tests, Logs, Error logging
