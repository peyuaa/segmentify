package handlers

import (
	"context"
	"net/http"

	"github.com/peyuaa/segmentify/data"
	"github.com/peyuaa/segmentify/models"
)

func (s *Segments) MiddlewareValidateSegment(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		segment := models.CreateSegmentRequest{}

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

func (s *Segments) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := models.UserSegments{}

		err := data.FromJSON(&user, r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			err = data.ToJSON(&GenericError{Message: err.Error()}, rw)
			if err != nil {
				s.l.Error("Unable to serialize GenericError", "error", err)
			}
		}

		errs := s.v.Validate(user)
		if len(errs) != 0 {
			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			err = data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			if err != nil {
				s.l.Error("Unable to serialize ValidationError", "error", err)
			}
			return
		}

		// add the request object to the context
		ctx := context.WithValue(r.Context(), KeyUserSegments{}, user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
