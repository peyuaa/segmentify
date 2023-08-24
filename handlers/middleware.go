package handlers

import (
	"context"
	"net/http"

	"github.com/peyuaa/segmentify/data"
)

func (s *Slugs) MiddlewareValidateSlug(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		slug := data.Slug{}

		err := data.FromJSON(&slug, r.Body)
		if err != nil {
			s.l.Error("Unable to deserialize slug", "error", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
		}

		// validate the slug
		errs := s.v.Validate(slug)
		if len(errs) != 0 {
			s.l.Error("Unable to validate slug", "error", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add the slug to the context
		ctx := context.WithValue(r.Context(), KeySlug{}, slug)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
