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

type Subscription struct {
	RepoName string   `json:"repoName"`
	ApiURL   string   `json:"apiUrl"`
	HtmlURL  string   `json:"htmlUrl"`
	Labels   []string `json:"labels"`
}

type LabelsByRepo map[string][]string
type VisitedRepositories map[string]bool

func GetSubscriptionsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	subscriptions := make([]Subscription, 0)

	data, err := models.GetSubscriptionsByUserID(userID)
	// If no repository found with the given name return an empty response
	if err == sql.ErrNoRows {
		utils.RespondWithJSON(w, http.StatusOK, subscriptions)
		return
	}

	labelsByRepo := LabelsByRepo{}
	visitedRepositories := VisitedRepositories{}
	for _, d := range data {
		labelsByRepo[d["repoName"]] = append(labelsByRepo[d["repoName"]], d["label"])
		visitedRepositories[d["repoName"]] = false
	}

	for _, d := range data {
		var subscription Subscription

		if !visitedRepositories[d["repoName"]] {
			visitedRepositories[d["repoName"]] = true

			subscription.RepoName = d["repoName"]
			subscription.HtmlURL = d["htmlURL"]
			subscription.Labels = labelsByRepo[d["repoName"]]

			subscriptions = append(subscriptions, subscription)
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, subscriptions)

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
	var subscription Subscription

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var repoID string
	repoID, err := models.GetRepositoryIDByName(subscription.RepoName)

	// If no repository found with the given name, create the repository in DB
	if err == sql.ErrNoRows {
		repoID, err = models.CreateRepository(
			subscription.RepoName,
			subscription.ApiURL,
			subscription.HtmlURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If no other err returned then save subscription and respond back
	err = models.CreateSubscription(userID, repoID, subscription.Labels)
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

func RemoveSubscription(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	var subscription Subscription

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var repoID string
	repoID, err := models.GetRepositoryIDByName(subscription.RepoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = models.RemoveSubscription(userID, repoID, subscription.Labels)
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
