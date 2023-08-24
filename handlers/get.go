package handlers

import (
	"net/http"

	"github.com/peyuaa/segmentify/data"
)

func (s *Slugs) Get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	slugs := data.GetSlugs()

	err := data.ToJSON(slugs, rw)
	if err != nil {
		s.l.Error("Unable to marshal json", "error", err)
	}
}
