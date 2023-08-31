package handlers

import (
	"github.com/peyuaa/segmentify/data"

	"github.com/charmbracelet/log"
)

// Segments is a struct that defines the handlers for the segments
type Segments struct {
	l *log.Logger
	v *data.Validation
	d *data.SegmentifyDB
}

// NewSegments returns a new Segments struct
func NewSegments(l *log.Logger, v *data.Validation, d *data.SegmentifyDB) *Segments {
	return &Segments{
		l: l,
		v: v,
		d: d,
	}
}

// KeySegment is a key used for the Segment object in the context
type KeySegment struct{}

// KeyUserSegments is a key used for UserSegments object in the context
type KeyUserSegments struct{}
