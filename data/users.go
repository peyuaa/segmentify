package data

import (
	"context"
	"fmt"
)

// SegmentAdd defines the structure for an API for adding segments
type SegmentAdd struct {
	// the segment's slug
	Slug string `json:"slug" validate:"required,min=5,max=50"`
	// expiration date
	Expired string `json:"omitempty,expired" validate:"omitempty,datetime=2006-01-02T15:04:05Z"`
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

func (s *SegmentifyDB) ChangeUserSegments(ctx context.Context, us UserSegments) error {
	// check if the add segments exists
	for _, segment := range us.AddSegments {
		_, err := s.GetSegmentBySlug(ctx, segment.Slug)
		if err != nil {
			return fmt.Errorf("unable to get segment \"%v\": %w", segment.Slug, err)
		}
	}

	// check if the remove segments exists
	for _, segment := range us.RemoveSegments {
		_, err := s.GetSegmentBySlug(ctx, segment.Slug)
		if err != nil {
			return fmt.Errorf("unable to get segment \"%v\": %w", segment.Slug, err)
		}
	}

	// add the segments to the user
	err := s.changeUsersSegments(ctx, us)
	if err != nil {
		return fmt.Errorf("unable to change user segments: %w", err)
	}

	return nil
}
