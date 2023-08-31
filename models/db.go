package models

import (
	"database/sql"
	"time"
)

// SegmentDB defines the structure for a segment in the database
type SegmentDB struct {
	// segment's id
	ID int

	// segment's slug
	Slug string

	// is the segment deleted
	IsDeleted bool
}

// SegmentsDB defines a slice of SegmentDB
type SegmentsDB []SegmentDB

// SegmentAddDB defines the structure for adding a segment to the database
type SegmentAddDB struct {
	// the segment's slug
	Slug string

	// expiration date
	Expired sql.NullString
}

// SegmentDeleteDB defines the structure for deleting a segment from the database
type SegmentDeleteDB struct {
	// the segment's slug
	Slug string
}

// UserSegmentsDB defines the structure for adding and removing segments from the user
type UserSegmentsDB struct {
	// user's id
	ID int

	// add the segments to the user
	AddSegments []SegmentAddDB

	// remove the segments from the user
	RemoveSegments []SegmentDeleteDB
}

// UserSegmentHistoryDB defines the structure for a segment in the database
type UserSegmentHistoryDB struct {
	// user's id
	ID int

	// segment's slug
	Slug string

	// date added
	DateAdded time.Time

	// date removed
	DateRemoved sql.NullTime
}

// UserSegmentsHistoryDB defines a slice of UserSegmentHistoryDB
type UserSegmentsHistoryDB []UserSegmentHistoryDB
