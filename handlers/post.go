package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/peyuaa/segmentify/data"
)

func (s *Segments) Create(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	// fetch the segment from the context
	segment := r.Context().Value(KeySegment{}).(data.Segment)

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
	segment, err = s.d.GetSegmentBySlug(r.Context(), segment.Slug)

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
		Path:   fmt.Sprintf("/segments/%s", segment.Slug),
	}
	rw.Header().Add("Location", u.String())

	rw.WriteHeader(http.StatusCreated)
	err = data.ToJSON(segment, rw)
	if err != nil {
		s.l.Error("Unable to serialize segment", "error", err)
	}
}
