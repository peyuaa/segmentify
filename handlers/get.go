package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/peyuaa/segmentify/data"

	"github.com/gorilla/mux"
)

const (
	MaxUserID = 2147483647
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

func (s *Segments) GetActiveSegments(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id, err := s.getUserId(r)
	if err != nil {
		s.writeGenericError(rw, http.StatusBadRequest, "", err)
		return
	}

	segments, err := s.d.GetUsersSegments(r.Context(), id)
	if err != nil {
		if errors.Is(err, data.ErrNoUserData) {
			s.writeGenericError(rw, http.StatusNotFound, "userID="+strconv.Itoa(id), err)
			return
		}
		s.writeInternalServerError(rw, "unable to get user's segments", err)
		return
	}

	err = data.ToJSON(segments, rw)
	if err != nil {
		s.writeInternalServerError(rw, "unable to marshal json", err)
	}
}

// getSlug returns the slug from the url
func (s *Segments) getSlug(r *http.Request) string {
	vars := mux.Vars(r)

	return vars["slug"]
}

func (s *Segments) getUserId(r *http.Request) (int, error) {
	id := mux.Vars(r)["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("unable to convert userID to int: %w", err)
	}

	if userID > MaxUserID {
		return 0, fmt.Errorf("userID is too big, got=%v, max value is %v", userID, MaxUserID)
	}

	return userID, nil
}
