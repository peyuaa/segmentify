package data

import "fmt"

// SegmentNotFound is an error raised when a segment can not be found in the database
var SegmentNotFound = fmt.Errorf("segment not found")

// Segment defines the structure for an API segment
type Segment struct {
	// the id for the segment
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the segment

	// the segment's slug
	//
	// required: true
	// max length: 255
	Slug string `json:"slug" validate:"required"`
}

var segments = []Segment{
	{
		ID:   1,
		Slug: "AVITO_VOICE_MESSAGES",
	},
	{
		ID:   2,
		Slug: "AVITO_PERFORMANCE_VAS",
	},
	{
		ID:   3,
		Slug: "AVITO_DISCOUNT_30",
	},
}

func (s *Segments) Add(ctx context.Context, segment Segment) error {
	exists, err := s.db.IsSegmentExists(ctx, segment.Slug)
	if err != nil {
		return fmt.Errorf("unable to check segment existence: %w", err)
	}
	if exists {
		return ErrSegmentAlreadyExists
	}

	err = s.db.InsertSegment(ctx, segment.Slug)
	if err != nil {
		return fmt.Errorf("unable to insert segment: %w", err)
	}

	return nil
}

func GetSegments() []Segment {
	return segments
}

func GetSegmentByID(id int) (*Segment, error) {
	for _, segment := range segments {
		if segment.ID == id {
			return &segment, nil
		}
	}

	return nil, ErrSegmentNotFound
}
