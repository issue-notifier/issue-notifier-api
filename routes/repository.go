package routes

import (
	"net/http"

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
