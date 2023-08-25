package handlers

import (
	"errors"
	"net/http"

	"github.com/peyuaa/segmentify/data"
)

func (s *Segments) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the segment from the context
	segment := r.Context().Value(KeySegment{}).(data.Segment)

	s.l.Debug("Inserting segment", "segment", segment)

	err := s.d.Add(r.Context(), segment)

	switch {
	case err == nil:
		rw.WriteHeader(http.StatusCreated)
	case errors.Is(err, data.ErrSegmentAlreadyExists):
		rw.WriteHeader(http.StatusConflict)
		err := data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to serialize GenericError", "error", err)
		}
	default:
		s.l.Error("Unable to insert segment", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
