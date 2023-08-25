package handlers

import (
	"fmt"
	"net/http"

	"github.com/peyuaa/segmentify/data"
)

var (
	// ErrInternalServer is an error message returned when something goes wrong
	// that do not depend on the user input
	// We don't want to expose the details of the error to the user, instead we log it
	ErrInternalServer = &GenericError{Message: "Dont worry, we are working on it!"}
)

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

func (s *Segments) writeGenericError(rw http.ResponseWriter, status int, message string, err error) {
	rw.WriteHeader(status)
	// specify the segment that was not found
	err = fmt.Errorf(message, "error", err)
	err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
	if err != nil {
		s.l.Error("Unable to serialize GenericError", "error", err)
	}
}

func (s *Segments) writeInternalServerError(rw http.ResponseWriter, logMessage string, err error) {
	s.l.Error(logMessage, "error", err)
	rw.WriteHeader(http.StatusInternalServerError)
	err = data.ToJSON(ErrInternalServer, rw)
	if err != nil {
		s.l.Error("Unable to serialize GenericError", "error", err)
	}
}
