package utils

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON utitlity function to convert HTTP response data to JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
