package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/peyuaa/segmentify/data"
	"github.com/peyuaa/segmentify/models"

	"github.com/gorilla/mux"
)

const (
	MaxUserID = 2147483647
)

// swagger:route GET /segments segments listSegments
// Returns a list of all segments from the database, deleted segments are included
//
// Produces:
// - application/json
//
// Schemes: http
//
// Responses:
// 	200: segmentsResponse
// 	500: errorResponse

// GetSegments returns the active segments from the database
func (s *Segments) GetSegments(rw http.ResponseWriter, r *http.Request) {
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

// swagger:route GET /segments/{Slug} segments getSegmentBySlug
// Returns a segment from the database by slug
//
// Produces:
// - application/json
//
// Schemes: http
//
// Parameters:
// 	+ name: Slug
// 	  in: path
// 	  description: slug of the segment
// 	  required: true
// 	  type: string
//
// Responses:
// 	200: segmentResponse
// 	404: errorResponse
// 	500: errorResponse

// GetBySlug returns a segment from the database by slug
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

// swagger:route GET /segments/users/{id} segments getActiveSegmentsForUser
// Returns a list of active segments for the user
//
// Produces:
// - application/json
//
// Schemes: http
//
// Responses:
// 	200: segmentsResponse
//	400: errorResponse
// 	404: errorResponse
// 	500: errorResponse

// GetActiveSegments returns the active segments for the user
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

// swagger:route GET /segments/users/{id}/history segments getUserHistory
// Returns a link to the user's segments history for the specified period
//
// Produces:
// - application/json
//
// Schemes: http
//
// Parameters:
// 	+ name: id
// 	  in: path
// 	  description: user id
// 	  required: true
// 	  type: integer
// 	+ name: from
// 	  in: query
// 	  description: start of the period. Format: YYYY-MM-DD
//	  required: true
// 	  type: string
// 	+ name: to
// 	  in: query
// 	  description: end of the period. Format: YYYY-MM-DD
// 	  required: true
// 	  type: string
//
// Responses:
// 	200: userHistoryResponse
// 	400: errorResponse
//	404: errorResponse
// 	500: errorResponse

// UserHistory returns the user's segments history for the specified period
func (s *Segments) UserHistory(rw http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserId(r)
	if err != nil {
		s.writeGenericError(rw, http.StatusBadRequest, "", err)
		return
	}

	from, to, err := s.getFromTo(r)
	if err != nil {
		s.writeGenericError(rw, http.StatusBadRequest, "unable to parse time", err)
		return
	}

	file, err := s.d.GetUserHistory(r.Context(), userID, from, to)
	if err != nil {
		if errors.Is(err, data.ErrNoUserHistoryData) {
			s.writeGenericError(rw, http.StatusNotFound, "userID="+strconv.Itoa(userID), err)
			return
		}
		s.writeInternalServerError(rw, "unable to get user's segments history", err)
		return
	}

	// create url for the link
	u := &url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   file,
	}

	history := models.UserHistoryResponse{
		Link: u.String(),
	}

	err = data.ToJSON(history, rw)
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

func (s *Segments) getFromTo(r *http.Request) (from, to time.Time, err error) {
	fromStr := r.URL.Query().Get("from")
	if fromStr == "" {
		return from, to, fmt.Errorf("from is empty")
	}

	from, err = time.Parse(time.DateOnly, fromStr)
	if err != nil {
		return from, to, fmt.Errorf("unable to parse from: %w", err)
	}

	toStr := r.URL.Query().Get("to")
	if toStr == "" {
		return from, to, fmt.Errorf("to is empty")
	}

	to, err = time.Parse(time.DateOnly, toStr)
	if err != nil {
		return from, to, fmt.Errorf("unable to parse to: %w", err)
	}

	// set to time to the end of the day
	to = to.Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59))

	// check if from is before to
	if from.After(to) {
		return from, to, fmt.Errorf("from is after to")
	}

	return from, to, nil
}
