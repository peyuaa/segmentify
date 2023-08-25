package handlers

import (
	"github.com/peyuaa/segmentify/data"

	"github.com/charmbracelet/log"
)

type Segments struct {
	l *log.Logger
	v *data.Validation
	d *data.SegmentifyDB
}

func NewSegments(l *log.Logger, v *data.Validation, d *data.SegmentifyDB) *Segments {
	return &Segments{
		l: l,
		v: v,
		d: d,
	}
}

// KeySegment is a key used for the Segment object in the context
type KeySegment struct{}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
