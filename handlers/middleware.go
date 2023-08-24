package handlers

import (
	"context"
	"net/http"

	"github.com/peyuaa/segmentify/data"
)

func (s *Segments) MiddlewareValidateSegment(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		segment := data.Segment{}

		err := data.FromJSON(&segment, r.Body)
		if err != nil {
			s.l.Error("Unable to deserialize segment", "error", err)

			rw.WriteHeader(http.StatusBadRequest)
			err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
			if err != nil {
				s.l.Error("Unable to serialize GenericError", "error", err)
			}
		}

		// validate the segment
		errs := s.v.Validate(segment)
		if len(errs) != 0 {
			s.l.Error("Unable to validate segment", "error", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			err = data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			if err != nil {
				s.l.Error("Unable to serialize ValidationError", "error", err)
			}
			return
		}

		// add the segment to the context
		ctx := context.WithValue(r.Context(), KeySegment{}, segment)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
