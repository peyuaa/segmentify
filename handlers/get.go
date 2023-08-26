package handlers

import (
	"errors"
	"net/http"

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

func (s *Segments) GetBySlug(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	slug := s.getSlug(r)

	segment, err := s.d.GetSegmentBySlug(r.Context(), slug)

	switch {
	case err == nil:
	case errors.Is(err, data.ErrSegmentNotFound):
		rw.WriteHeader(http.StatusNotFound)
		err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to marshal json", "error", err)
		}
		return
	default:
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

// getSlug returns the slug from the url
func (s *Segments) getSlug(r *http.Request) string {
	// parse the product id from the url
	vars := mux.Vars(r)

	return vars["slug"]
}
