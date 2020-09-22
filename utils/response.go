package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Payload interface{} `json:"payload"`
}

func CreateResponse(success bool, message string, payload interface{}) (response *Response) {
	response = &Response{
		Success: success,
		Error:   message,
		Payload: payload,
	}

	return
}

func RespondWithJSON(w http.ResponseWriter, code int, response *Response) {
	res, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
