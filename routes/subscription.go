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

// GetSubscriptionsByUserID godoc
// @Summary Get all subscriptions for the authenticated user
// @Description Get all subscriptions for the authenticated user
// @Tags subscription
// @Security Github OAuth
// @Produce json
// @Success 200 {array} models.Subscription
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/user/subscription/view [get]
func GetSubscriptionsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	subscriptions, err := models.GetSubscriptionsByUserID(userID)
	// If no repository found with the given name return an empty response
	if err != nil && err != sql.ErrNoRows {
		utils.LogError.Println("Failed to fetch subscriptions for userID:", userID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogInfo.Println("Returning", len(subscriptions), "subscription for userID:", userID)
	utils.RespondWithJSON(w, http.StatusOK, subscriptions)
}

// GetSubscriptionsByRepoID godoc
// @Summary Get all subscriptions for the given `repoID`
// @Description Get all subscriptions for the given `repoID`
// @Tags subscription
// @Produce json
// @Param repoID path string true "Repository ID for which subscription data needs to be fetched"
// @Success 200 {array} models.UserIDLabel
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/{repoID}/subscription/view [get]
func GetSubscriptionsByRepoID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoID := vars["repoID"]

	var data []models.UserIDLabel
	data, err := models.GetSubscriptionsByRepoID(repoID)
	if err != nil {
		utils.LogError.Println("Failed to subscriptions for repoID:", repoID, ". Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogInfo.Println("Returning", len(data), "subscription for repoID:", repoID)
	utils.RespondWithJSON(w, http.StatusOK, data)
}

// GetSubscribedLabelsByUserIDAndRepoName godoc
// @Summary Get all subscriptions for the given `userID` and `repoName` of the authenticated user
// @Description Get all subscriptions for the given `userID` and `repoName` of the authenticated user
// @Tags subscription
// @Produce json
// @Security Github OAuth
// @Param repoName path string true "Repository Name for which subscription data needs to be fetched. Format `facebook/react`"
// @Success 200 {array} models.Labels
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/user/{repoName}/subscription/labels [get]
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
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, labels)
}

// CreateSubscriptions godoc
// @Summary Create new subscription for the authenticated user
// @Description Create new subscription for the authenticated user
// @Tags subscription
// @Produce json
// @Security Github OAuth
// @Param subscriptions body models.Subscription true "Repository Name and list of Labels to create a new subscription"
// @Success 201 {string} Create Success
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/user/subscription/add [post]
func CreateSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	var subscription models.Subscription

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

	utils.RespondWithJSON(w, http.StatusCreated, "Create Success")
}

// UpdateSubscriptions godoc
// @Summary Update existing subscription for the authenticated user
// @Description Update existing subscription for the authenticated user
// @Tags subscription
// @Produce json
// @Security Github OAuth
// @Param subscriptions body models.Subscription true "Repository Name and list of Labels which needs to be added to the existing subscription"
// @Success 204 {string} Update Success
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/user/subscription/update [put]
func UpdateSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	var subscription models.Subscription

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

	utils.RespondWithJSON(w, http.StatusNoContent, "Update Success")

}

// RemoveSubscriptions godoc
// @Summary Deletes existing subscription for the authenticated user
// @Description Deletes existing subscription for the authenticated user
// @Tags subscription
// @Produce json
// @Security Github OAuth
// @Param subscriptions body models.Subscription true "Repository Name and list of Labels which needs to be deleted from the existing subscription"
// @Success 204 {string} Remove Success
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/user/subscription/remove [delete]
func RemoveSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(w, r)

	var subscription models.Subscription

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

	utils.RespondWithJSON(w, http.StatusNoContent, "Remove Success")
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
