package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/peyuaa/segmentify/data"

	"github.com/gorilla/mux"
)

func (s *Segments) Get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	segments, err := s.d.GetSegments(r.Context())
	if err != nil {
		s.l.Error("Unable to get segments", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to marshal json", "error", err)
		}
		return
	}

	err = data.ToJSON(segments, rw)
	if err != nil {
		s.l.Error("Unable to marshal json", "error", err)
	}
}

func (s *Segments) GetById(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := s.getId(r)

	segment, err := s.d.GetSegmentByID(r.Context(), id)

	switch {
	case err == nil:
	case errors.Is(err, data.ErrSegmentNotFound):
		s.l.Warn("Unable to find segment in database", "id", id, "error", err)
		rw.WriteHeader(http.StatusNotFound)
		err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to marshal json", "error", err)
		}
		return
	default:
		s.l.Error("Error retrieving segment from the database", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to marshal json", "error", err)
		}
		return
	}

	err = data.ToJSON(segment, rw)
	if err != nil {
		s.l.Error("Unable to marshal json", "error", err)
	}
}

// getId returns the slug id from the url
// Log error if func cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func (s *Segments) getId(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.l.Error("[SHOULD NEVER HAPPEN] Unable to convert id into integer", "error", err)
	}

	return id
}
