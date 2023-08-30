package handlers

import (
	"errors"
	"net/http"

	"github.com/peyuaa/segmentify/data"
)

// swagger:route DELETE /segments/{Slug} segments deleteSegment
// Deletes a segment from the database
// responses:
// 	204: noContentResponse
// 	404: errorResponse
// 	500: errorResponse

// Delete handles DELETE requests and mark segment as deleted in the database
// The reason why we don't delete the segment from the database is that we want to keep
// the history of the already used segments
func (s *Segments) Delete(rw http.ResponseWriter, r *http.Request) {
	slug := s.getSlug(r)

	err := s.d.Delete(r.Context(), slug)
	switch {
	case err == nil:
		rw.WriteHeader(http.StatusNoContent)
	case errors.Is(err, data.ErrSegmentNotFound):
		rw.WriteHeader(http.StatusNotFound)
		err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to serialize GenericError", "error", err)
		}
	default:
		s.l.Error("Unable to delete segment", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
		if err != nil {
			s.l.Error("Unable to serialize GenericError", "error", err)
		}
	}
}
