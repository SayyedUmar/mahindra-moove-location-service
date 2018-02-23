package web

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Error struct {
	Message string `json:"error"`
}

func ErrorWithMessage(msg string) *Error {
	return &Error{Message: msg}
}

func (e *Error) Respond(w http.ResponseWriter, statusCode int) {
	log.Error(e.Message)
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(e)
	if err != nil {
		log.Warn(err)
	}
}

func writeOk(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(map[string]string{"message": "success"})
	if err != nil {
		ErrorWithMessage(err.Error()).Respond(w, 200)
	}
}
