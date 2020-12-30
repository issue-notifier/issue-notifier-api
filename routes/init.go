package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/issue-notifier/issue-notifier-api/middleware"
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
func Init(githubClientID, githubClientSecret string) {
	GITHUB_CLIENT_ID = githubClientID
	GITHUB_CLIENT_SECRET = githubClientSecret

	Router = mux.NewRouter()
	noAuthRouter = Router.PathPrefix("/api/v1").Subrouter()
	authRouter = noAuthRouter.PathPrefix("/user").Subrouter()
	authRouter.Use(middleware.IsAuthenticated)

	Router.HandleFunc("/health", health).Methods("GET")

	// matches /api/v1/...
	noAuthRouter.HandleFunc("/login/github/oauth2", GitHubLogin).Methods("GET")
	noAuthRouter.HandleFunc("/repositories", GetAllRepositories).Methods("GET")
	noAuthRouter.HandleFunc("/repository/{repoID}/update/lastEventAt", UpdateLastEventAtByRepoID).Methods("PUT")
	noAuthRouter.HandleFunc("/subscription/{repoID}/view", GetSubscriptionsByRepoID).Methods("GET")

	// matches /api/v1/user/...
	authRouter.HandleFunc("/authenticated", GetAuthenticatedUser).Methods("GET")
	authRouter.HandleFunc("/logout", Logout).Methods("GET")
	authRouter.HandleFunc("/subscription/add", CreateSubscriptions).Methods("POST")
	authRouter.HandleFunc("/subscription/view", GetSubscriptionsByUserID).Methods("GET")
	authRouter.HandleFunc("/subscription/update", UpdateSubscriptions).Methods("PUT")
	authRouter.HandleFunc("/subscription/remove", RemoveSubscriptions).Methods("DELETE")
	authRouter.HandleFunc("/subscription/{org}/{repo}/labels", GetSubscribedLabelsByUserIDAndRepoName).Methods("GET")
}

func health(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(map[string]string{
		"status": "UP",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
