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

// Init all types of routes
func Init() {
	Router = mux.NewRouter()
	noAuthRouter = Router.PathPrefix("/api/v1").Subrouter()
	authRouter = noAuthRouter.PathPrefix("/user").Subrouter()
	authRouter.Use(middleware.IsAuthenticated)

	Router.HandleFunc("/health", healthCheck).Methods("GET")

	// matches /api/v1/...
	noAuthRouter.HandleFunc("/login/github/oauth2", GitHubLogin).Methods("GET")

	// matches /api/v1/user/...
	authRouter.HandleFunc("/authenticated", GetAuthenticatedUser).Methods("GET")
	authRouter.HandleFunc("/logout", Logout).Methods("GET")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(map[string]string{
		"status": "UP",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
