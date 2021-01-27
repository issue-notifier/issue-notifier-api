package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/issue-notifier/issue-notifier-api/models"
	"github.com/issue-notifier/issue-notifier-api/utils"
)

// GetAllRepositories godoc
// @Summary Get all repositories
// @Description Get complete repository information for all repositories from the database
// @Tags repository
// @Produce json
// @Success 200 {array} models.Repository
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/repositories [get]
func GetAllRepositories(w http.ResponseWriter, r *http.Request) {
	repositories, err := models.GetAllRepositories()
	if err != nil {
		utils.LogError.Println("Failed to fetch repositories data from the database. Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogInfo.Println("Returning", len(repositories), "repositories from the database")
	utils.RespondWithJSON(w, http.StatusOK, repositories)
}

type lastEventAt struct {
	LastEventAt time.Time `json:"lastEventAt" db:"last_event_at"`
}

var layout string = "2006-01-02 15:04:05-07:00"

// UpdateLastEventAtByRepoID godoc
// @Summary Updates `lastEventAt` time for the given `repoID`
// @Description Updates `lastEventAt` time for the given `repoID` in the database
// @Tags repository
// @Consume json
// @Produce json
// @Param repoID path string true "Repository ID for which `lastEventAt` needs to be updated"
// @Param lastEventAt body lastEventAt true "`lastEventAt` time (with timezone) at which last event for the repository occurred. Format `lastEventAt`: `2006-01-02 15:04:05-07:00`"
// @Success 204 {string} Updated LastEventAt Successfully "No Content"
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/repository/{repoID}/update/lastEventAt [put]
func UpdateLastEventAtByRepoID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoID := vars["repoID"]

	var data lastEventAt

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		utils.LogError.Println("Failed to decode the request body. Error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := models.UpdateLastEventAtByRepoID(repoID, data.LastEventAt.Format(layout))
	if err != nil {
		utils.LogError.Println("Failed to updated `lastEventAt` time for repoID:", repoID, ". Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogInfo.Println("Successfully updated `lastEventAt` time for repoID:", repoID)
	utils.RespondWithJSON(w, http.StatusNoContent, "Updated LastEventAt Successfully")
}
