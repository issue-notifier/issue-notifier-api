package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/issue-notifier/issue-notifier-api/models"
	"github.com/issue-notifier/issue-notifier-api/session"
	"github.com/issue-notifier/issue-notifier-api/utils"
	"github.com/lib/pq"
)

type Subscription struct {
	RepoName string        `json:"repoName"`
	Labels   models.Labels `json:"labels"`
}

func GetSubscriptionsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	subscriptions, err := models.GetSubscriptionsByUserID(userID)
	// If no repository found with the given name return an empty response
	if err == sql.ErrNoRows {
		utils.RespondWithJSON(w, http.StatusOK, subscriptions)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, subscriptions)
}

func GetSubscriptionsByRepoID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoID := vars["repoID"]

	var data []map[string]interface{}

	data, err := models.GetSubscriptionsByRepoID(repoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, data)
}

func GetSubscribedLabelsByUserIDAndRepoName(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	vars := mux.Vars(r)
	repoName := vars["org"] + "/" + vars["repo"]

	var labels models.Labels

	// TODO: Convert this into one single DB call using JOINTS
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

func CreateSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)
	var subscription Subscription

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	repoID, err := models.GetRepositoryIDByName(subscription.RepoName)

	// If no repository found with the given name, create the repository in DB
	if err == sql.ErrNoRows {
		repoID, err = models.CreateRepository(subscription.RepoName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save subscription and respond back
	err = models.CreateSubscriptions(userID, repoID, subscription.Labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Create Success")
}

func UpdateSubscriptions(w http.ResponseWriter, r *http.Request) {

	userID := getUserIDFromSession(w, r)
	var subscription Subscription

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	repoID, err := models.GetRepositoryIDByName(subscription.RepoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingLabels, err := models.GetSubscribedLabelsByUserIDAndRepoID(userID, repoID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Make it more efficient if possible
	// Check if the new labels already exist or not in the database
	// Or create a map and resave the same label and don't throw any error
	for _, existingLabel := range existingLabels {
		for _, newLabel := range subscription.Labels {
			if existingLabel.Name == newLabel.Name {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	// Save subscription and respond back
	err = models.UpdateSubscriptions(userID, repoID, subscription.Labels)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			http.Error(w, "You have already subscribed to these labels", http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Update Success")

}

func RemoveSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	var subscription Subscription

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	repoID, err := models.GetRepositoryIDByName(subscription.RepoName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = models.RemoveSubscriptions(userID, repoID, subscription.Labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = models.DeleteRepositoriesWithNoLabels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Remove Success")
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
