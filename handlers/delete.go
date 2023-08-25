package handlers

import (
	"errors"
	"net/http"

	"github.com/peyuaa/segmentify/data"
)

func (s *Segments) Delete(rw http.ResponseWriter, r *http.Request) {
	slug := s.getSlug(r)

	err := s.d.Delete(r.Context(), slug)
	switch {
	case err == nil:
		rw.WriteHeader(http.StatusNoContent)
	case errors.Is(err, data.ErrSegmentNotFound):
		rw.WriteHeader(http.StatusNotFound)
		err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
	default:
		s.l.Error("Unable to delete segment", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
	}
}
