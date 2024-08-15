package response

import (
	"net/http"
)

func SendError(w http.ResponseWriter, status int, message interface{}) {
	err := WriteJSON(w, status, Envelope{"error": message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SendServerError(w http.ResponseWriter) {
	message := "the server encountered a problem and could not process your request"
	SendError(w, http.StatusInternalServerError, message)
}
