package routes

import (
	"encoding/json"
	"net/http"

	"github.com/issue-notifier/issue-notifier-api/middleware"
	"github.com/issue-notifier/issue-notifier-api/utils"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router main Gorilla mux router of app
var (
	Router       *mux.Router
	noAuthRouter *mux.Router
	authRouter   *mux.Router
)

// Github credentials
var (
	GithubClientID     string
	GithubClientSecret string
)

// Init initializes all types of routes
func Init(githubClientID, githubClientSecret string) {
	GithubClientID = githubClientID
	GithubClientSecret = githubClientSecret

	Router = mux.NewRouter()
	Router.Use(middleware.LogHTTPRequest)
	noAuthRouter = Router.PathPrefix("/api/v1").Subrouter()
	authRouter = noAuthRouter.PathPrefix("/user").Subrouter()
	authRouter.Use(middleware.IsAuthenticated)

	// matches /health
	Router.HandleFunc("/health", health).Methods("GET")

	noAuthRouter.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

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

	utils.LogInfo.Println("Successfully initialized routes")
}

func health(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(map[string]string{
		"status": "UP",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
