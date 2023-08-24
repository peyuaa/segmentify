package handlers

import (
	"net/http"

	"github.com/peyuaa/segmentify/data"
)

func (s *Segments) Create(_ http.ResponseWriter, r *http.Request) {
	// fetch the segment from the context
	segment := r.Context().Value(KeySegment{}).(data.Segment)

	s.l.Debug("Inserting segment", "segment", segment)
	data.AddSegment(segment)
}
