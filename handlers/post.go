package handlers

import (
	"net/http"

	"github.com/peyuaa/segmentify/data"
)

func (s *Slugs) Create(_ http.ResponseWriter, r *http.Request) {
	// fetch the slug from the context
	slug := r.Context().Value(KeySlug{}).(data.Segment)

	s.l.Debug("Inserting slug", "slug", slug)
	data.AddSegment(slug)
}
