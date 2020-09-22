package main

import (
	"log"
	"net/http"

	"github.com/hemakshis/issue-notifier-api/database"
	"github.com/hemakshis/issue-notifier-api/routes"
	"github.com/hemakshis/issue-notifier-api/session"
)

func main() {
	database.Init()

	routes.Init()

	session.Init()

	log.Fatal(http.ListenAndServe(":8001", routes.Router)) // TODO: Move port also to env vars
}
