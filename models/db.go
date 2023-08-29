package models

import (
	"database/sql"
	"time"
)

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

type SegmentAddDB struct {
	// the segment's slug
	Slug string

	// expiration date
	Expired sql.NullString
}

type SegmentDeleteDB struct {
	// the segment's slug
	Slug string
}

type UserSegmentsDB struct {
	// user's id
	ID int

	// add the segments to the user
	AddSegments []SegmentAddDB

	// remove the segments from the user
	RemoveSegments []SegmentDeleteDB
}

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

type UserSegmentsHistoryDB []UserSegmentHistoryDB
