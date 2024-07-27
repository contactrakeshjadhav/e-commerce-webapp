package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// JsonResponse generates a new HTTP Response in Json format,
// acording to the HTTP status code it receives
func JsonResponse(w http.ResponseWriter, data []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

// NewApplicationError generates a new ApplicationError type
func NewApplicationError(err error) ApplicationError {
	return ApplicationError{
		Error: Error{Message: err.Error()},
	}
}

// ApplicationErrorResponse prepare and send a new JSON object with the application error
func ApplicationErrorResponse(rw http.ResponseWriter, err error, statusCode int) {
	data := NewApplicationError(err)
	resp, _ := json.Marshal(data)
	JsonResponse(rw, resp, statusCode)
}

func Decode(reader io.Reader, destination interface{}) error {
	err := json.NewDecoder(reader).Decode(destination)

	if err != nil {
		return errors.New(CheckErrors(err))
	}
	return nil
}
