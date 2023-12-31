package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/peyuaa/segmentify/data"
	"github.com/peyuaa/segmentify/models"
)

// swagger:route POST /segments segments createSegment
// Creates a new segment in the database
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Schemes: http
//
// Parameters:
//	+ name: segment
// 	  in: body
// 	  description: slug of the segment
// 	  required: true
// 	  type: createSegmentRequest
//
// Responses:
// 	201: createSegmentResponse
// 	400: errorResponse
// 	409: errorResponse
// 	500: errorResponse

// CreateSegment creates a new segment in the database
func (s *Segments) CreateSegment(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	// fetch the segment from the context
	segment := r.Context().Value(KeySegment{}).(models.CreateSegmentRequest)

	s.l.Debug("Inserting segment", "segment", segment)

	err := s.d.Add(r.Context(), segment)

	switch {
	case err == nil:
	case errors.Is(err, data.ErrSegmentAlreadyExists):
		rw.WriteHeader(http.StatusConflict)
		err := data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to serialize GenericError", "error", err)
		}
		return
	default:
		s.l.Error("Unable to insert segment", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		err := data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to serialize GenericError", "error", err)
		}
		return
	}

	// retrieve segment to include the result of the operation in the response body
	createdSegment, err := s.d.GetSegmentBySlug(r.Context(), segment.Slug)

	switch {
	case err == nil:
	case errors.Is(err, data.ErrSegmentNotFound):
		s.l.Error("Unable to find created segment", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		err := data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to serialize GenericError", "error", err)
		}
		return
	default:
		s.l.Error("Unable to find segment", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		err := data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to serialize GenericError", "error", err)
		}
		return
	}

	// set the Location header to the URL of the newly created resource
	u := &url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   fmt.Sprintf("/segments/%s", createdSegment.Slug),
	}
	rw.Header().Add("Location", u.String())

	rw.WriteHeader(http.StatusCreated)
	err = data.ToJSON(createdSegment, rw)
	if err != nil {
		s.l.Error("Unable to serialize segment", "error", err)
	}
}

// swagger:route POST /segments/users segments changeUsersSegments
// Add or remove segments from a user
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Schemes: http
//
// Parameters:
//	+ name: userSegments
// 	  in: body
// 	  description: segments to add and remove from the user
// 	  required: true
// 	  type: userSegmentsRequest
//
// Responses:
// 	200: activeSegmentsResponse
// 	400: errorResponse
// 	404: errorResponse
// 	500: errorResponse

// ChangeUsersSegments changes the segments of a user
// First it checks that all segments exist
// Then it adds the segments to the user
// Then it removes the segments from the user
// If add and delete segments contains the same segment, behavior is undefined. Don't do that.
func (s *Segments) ChangeUsersSegments(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the user segments from the context
	userSegments := r.Context().Value(KeyUserSegments{}).(models.UserSegmentsRequest)

	// add the segments to the user
	err := s.d.ChangeUserSegments(r.Context(), userSegments)

	switch {
	case err == nil:
	case errors.Is(err, data.ErrSegmentNotFound):
		s.writeGenericError(rw, http.StatusNotFound, "request contains unknown segments", err)
		return
	case errors.Is(err, data.ErrIncorrectChangeUserSegmentsRequest):
		s.writeGenericError(rw, http.StatusBadRequest, "request is incorrect", err)
		return
	case errors.Is(err, data.ErrSegmentDeleted):
		s.writeGenericError(rw, http.StatusBadRequest, "request contains deleted segment", err)
		return
	default:
		s.writeInternalServerError(rw, "Failed to change user segments", err)
		return
	}

	// get the user segments after the change
	newSegments, err := s.d.GetUsersSegments(r.Context(), userSegments.ID)
	switch {
	case err == nil:
	case errors.Is(err, data.ErrNoUserData):
	default:
		s.writeInternalServerError(rw, "unable to retrieve users segments", err)
	}

	response := models.ActiveSegmentsResponse{
		ActiveSegments: newSegments,
	}

	err = data.ToJSON(response, rw)
	if err != nil {
		s.writeInternalServerError(rw, "unable to serialize models.ActiveSegmentsResponse", err)
		return
	}
}
