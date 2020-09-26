package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hemakshis/issue-notifier-api/models"
	"github.com/hemakshis/issue-notifier-api/session"
	"github.com/hemakshis/issue-notifier-api/utils"
	"github.com/lib/pq"
)

type SubscriptionRequest struct {
	RepoName string   `json:"repoName"`
	ApiURL   string   `json:"apiUrl"`
	HtmlURL  string   `json:"htmlUrl"`
	Labels   []string `json:"labels"`
}

func GetSubscribedLabelsByUserIDAndRepoID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	vars := mux.Vars(r)
	repoName := vars["org"] + "/" + vars["repo"]

	labels := make([]string, 0)

	repoID, err := models.GetRepositoryIDByName(repoName)
	// If no repository found with the given name return an empty response
	if err == sql.ErrNoRows {
		utils.RespondWithJSON(w, http.StatusOK, labels)
		return
	}

	labels, err = models.GetSubscribedLabelsByUserIDAndRepoID(userID, repoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, labels)
}

func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)
	var subscriptionRequest SubscriptionRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&subscriptionRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var repoID string
	repoID, err := models.GetRepositoryIDByName(subscriptionRequest.RepoName)

	// If no repository found with the given name, create the repository in DB
	if err == sql.ErrNoRows {
		repoID, err = models.CreateRepository(
			subscriptionRequest.RepoName,
			subscriptionRequest.ApiURL,
			subscriptionRequest.HtmlURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If no other err returned then save subscription and respond back
	err = models.CreateSubscription(userID, repoID, subscriptionRequest.Labels)
	if err != nil && (err.(*pq.Error)).Code == "23505" {
		http.Error(w, "You have already subscribed to these labels", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Success")
}

func getUserIDFromSession(w http.ResponseWriter, r *http.Request) (userID string) {
	ses, err := session.Store.Get(r, session.CookieName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userSession, _ := ses.Values["UserSession"].(session.UserSession)
	userID = userSession.UserID

	return
}
