package models

import "time"

// Segment defines the structure for an API segment
type Segment struct {
	// the id for the segment
	ID int `json:"id"` // Unique identifier for the segment

	// the segment's slug
	Slug string `json:"slug" validate:"required,min=5,max=50"`

	// is the segment deleted
	IsDeleted bool `json:"is_deleted"`
}

// CreateSegmentRequest defines the structure for an API request for adding segments
// swagger:model createSegmentRequest
type CreateSegmentRequest struct {
	// the segment's slug
	//
	// required: true
	// min length: 5
	// max length: 50
	// example: AVITO_DISCOUNT_30
	Slug string `json:"slug" validate:"required,min=5,max=50"`
}

// ActiveSegment defines the structure of Segment for an API response for active user's segments
type ActiveSegment struct {
	// the segment's slug
	Slug string `json:"slug"`
}

// ActiveSegments defines the structure of response for active user's segments
type ActiveSegments []ActiveSegment

type ActiveSegmentsResponse struct {
	ActiveSegments ActiveSegments `json:"segments"`
}

// Segments defines a slice of Segment
type Segments []Segment

// SegmentAdd defines the structure for an API for adding segments
// swagger:model segmentAdd
type SegmentAdd struct {
	// the segment's slug
	//
	// required: true
	// min length: 5
	// max length: 50
	// example: AVITO_DISCOUNT_50
	Slug string `json:"slug" validate:"required,min=5,max=50"`

	// expiration date
	//
	// required: false
	// example: 2025-01-02T15:04:06Z
	Expired string `json:"expired,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z"`
}

// SegmentDelete defines the structure for an API for deleting segments
// swagger:model segmentDelete
type SegmentDelete struct {
	// the segment's slug
	//
	// required: true
	// min length: 5
	// max length: 50
	// example: AVITO_PERFORMANCE_VAS
	Slug string `json:"slug" validate:"required,min=5,max=50"`
}

// UserSegmentsRequest defines the structure for an API for adding segments to user
// swagger:model userSegmentsRequest
type UserSegmentsRequest struct {
	// user's id
	//
	// required: true
	// min: 1
	// max: 2147483647
	// example: 42
	ID int `json:"id" validate:"required,gt=0,number"`

	// add the segments to the user
	AddSegments []SegmentAdd `json:"add" validate:"dive"`

	// remove the segments from the user
	RemoveSegments []SegmentDelete `json:"remove" validate:"dive"`
}

// UserHistoryResponse defines the structure for an API response for getting user's segments history
type UserHistoryResponse struct {
	// link to csv file with user's segments history for specified period
	Link string `json:"link"`
}

// UserHistoryEntry defines user's segment history entry
type UserHistoryEntry struct {
	// userID
	ID int

	// segment's slug
	Slug string

	// operation type
	Operation string

	// date
	Date time.Time
}

// UserHistory defines a slice of UserHistoryEntry
// Implements sort.Interface
type UserHistory []UserHistoryEntry

// Len returns the length of UserHistory
func (u UserHistory) Len() int {
	return len(u)
}

// Less returns true if the date of the first entry is before the date of the second entry
func (u UserHistory) Less(i, j int) bool {
	return u[i].Date.Before(u[j].Date)
}

// Swap swaps the entries
func (u UserHistory) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
