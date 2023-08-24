package handlers

import (
	"github.com/peyuaa/segmentify/data"

	"github.com/charmbracelet/log"
)

type Slugs struct {
	l *log.Logger
	v *data.Validation
}

func NewSlugs(l *log.Logger, v *data.Validation) *Slugs {
	return &Slugs{
		l: l,
		v: v,
	}
}

// KeySlug is a key used for the Segment object in the context
type KeySlug struct{}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
