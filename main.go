package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/issue-notifier/issue-notifier-api/database"
	"github.com/issue-notifier/issue-notifier-api/routes"
	"github.com/issue-notifier/issue-notifier-api/session"
	"github.com/issue-notifier/issue-notifier-api/utils"

	_ "github.com/issue-notifier/issue-notifier-api/docs" // Generate Swagger Doc for APIs. Used
)

// Env vars
var (
	environment string
	port        string

	dbUser string
	dbPass string
	dbName string
	dbURL  string

	sessionAuthKey string

	githubClientID     string
	githubClientSecret string
)

// @title Github Issue-Notifier API
// @version 1.0
// @description APIs for the Github Issue Notifier Project. https://github.com/issue-notifier
// @termsOfService http://swagger.io/terms/
// @contact.name Hemakshi Sachdev
// @contact.email sachdev.hemakshi@gmail.com
// @host localhost:8001
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. Error:", err, "This is okay if the app is running on Production")
	}

	environment = os.Getenv("ENVIRONMENT")
	port = os.Getenv("PORT")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
	if environment == "production" {
		dbURL = os.Getenv("DATABASE_URL")
	} else {
		dbURL = ""
	}
	sessionAuthKey = os.Getenv("SESSION_AUTH_KEY")
	githubClientID = os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")

	utils.InitLogging(environment)

	database.Init(environment, dbUser, dbPass, dbName, dbURL)
	defer database.DB.Close()

	routes.Init(githubClientID, githubClientSecret)

	session.Init(sessionAuthKey)

	utils.LogError.Fatal(http.ListenAndServe(":"+port, routes.Router))
}
