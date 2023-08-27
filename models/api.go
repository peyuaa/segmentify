package models

// Segment defines the structure for an API segment
type Segment struct {
	// the id for the segment
	ID int `json:"id"` // Unique identifier for the segment

	// the segment's slug
	Slug string `json:"slug" validate:"required,min=5,max=50"`

	// is the segment deleted
	IsDeleted bool `json:"is_deleted"`
}

// ActiveSegment defines the structure of Segment for an API response for active user's segments
type ActiveSegment struct {
	// the segment's slug
	Slug string `json:"slug"`
}

// ActiveSegments defines the structure of response for active user's segments
type ActiveSegments []ActiveSegment

// Segments defines a slice of Segment
type Segments []Segment

// SegmentAdd defines the structure for an API for adding segments
type SegmentAdd struct {
	// the segment's slug
	Slug string `json:"slug" validate:"required,min=5,max=50"`
	// expiration date
	Expired string `json:"expired,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z"`
}

// SegmentDelete defines the structure for an API for deleting segments
type SegmentDelete struct {
	// the segment's slug
	Slug string `json:"slug" validate:"required,min=5,max=50"`
}

// UserSegments defines the structure for an API for adding segments to user
type UserSegments struct {
	// user's id
	ID int `json:"id" validate:"required,gt=0,number"`

	// add the segments to the user
	AddSegments []SegmentAdd `json:"add" validate:"dive"`

	// remove the segments from the user
	RemoveSegments []SegmentDelete `json:"remove" validate:"dive"`
}
