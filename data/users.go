package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/peyuaa/segmentify/models"
)

func (s *SegmentifyDB) ChangeUserSegments(ctx context.Context, us models.UserSegments) error {
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

	userSegmentsDB := models.UserSegmentsDB{
		ID:             us.ID,
		AddSegments:    make([]models.SegmentAddDB, len(us.AddSegments)),
		RemoveSegments: make([]models.SegmentDeleteDB, len(us.RemoveSegments)),
	}

	for i, segment := range us.AddSegments {
		userSegmentsDB.AddSegments[i] = models.SegmentAddDB{
			Slug: segment.Slug,
			Expired: sql.NullString{
				String: segment.Expired,
				Valid:  true,
			},
		}

		if segment.Expired == "" {
			userSegmentsDB.AddSegments[i].Expired.Valid = false
		}
	}

	for i, segment := range us.RemoveSegments {
		userSegmentsDB.RemoveSegments[i] = models.SegmentDeleteDB{
			Slug: segment.Slug,
		}
	}

	// add the segments to the user
	err := s.db.ChangeUsersSegments(ctx, userSegmentsDB)
	if err != nil {
		return fmt.Errorf("unable to change user segments: %w", err)
	}

	return nil
}
