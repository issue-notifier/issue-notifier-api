package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hemakshis/issue-notifier-api/middleware"
)

// Router main Gorilla mux router of app
var (
	Router       *mux.Router
	noAuthRouter *mux.Router
	authRouter   *mux.Router
)

var (
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
)

// Init all types of routes
func Init(githubClientId, githubClientSecret string) {
	GITHUB_CLIENT_ID = githubClientId
	GITHUB_CLIENT_SECRET = githubClientSecret

	Router = mux.NewRouter()
	noAuthRouter = Router.PathPrefix("/api/v1").Subrouter()
	authRouter = noAuthRouter.PathPrefix("/user").Subrouter()
	authRouter.Use(middleware.IsAuthenticated)

	Router.HandleFunc("/health", health).Methods("GET")

	// matches /api/v1/...
	noAuthRouter.HandleFunc("/login/github/oauth2", GitHubLogin).Methods("GET")

	// matches /api/v1/user/...
	authRouter.HandleFunc("/authenticated", GetAuthenticatedUser).Methods("GET")
	authRouter.HandleFunc("/logout", Logout).Methods("GET")
}

func health(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(map[string]string{
		"status": "UP",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
