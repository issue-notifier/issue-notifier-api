package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/issue-notifier/issue-notifier-api/models"
	"github.com/issue-notifier/issue-notifier-api/utils"
)

func GetAllRepositories(w http.ResponseWriter, r *http.Request) {
	repositories, err := models.GetAllRepositories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, repositories)
}

func UpdateLastEventAtByRepoID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoID := vars["repoID"]

	var data map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := models.UpdateLastEventAtByRepoID(repoID, data["lastEventAt"].(string))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Updated LastEventAt Successfully")
}
