package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/peyuaa/segmentify/models"
)

func (s *SegmentifyDB) ChangeUserSegments(ctx context.Context, us models.UserSegments) error {
	// check if the add segments exists
	for _, segment := range us.AddSegments {
		got, err := s.GetSegmentBySlug(ctx, segment.Slug)
		if err != nil {
			return fmt.Errorf("unable to get segment \"%v\": %w", segment.Slug, err)
		}
		if got.IsDeleted {
			return fmt.Errorf("can't add deleted segment \"%v\" to user: %w", segment.Slug, ErrSegmentDeleted)
		}
	}

	// check if the remove segments exists
	for _, segment := range us.RemoveSegments {
		_, err := s.GetSegmentBySlug(ctx, segment.Slug)
		if err != nil {
			return fmt.Errorf("unable to get segment \"%v\": %w", segment.Slug, err)
		}
	}

	// get user's segments
	userSegments, err := s.db.GetUsersSegments(ctx, us.ID)
	if err != nil {
		return fmt.Errorf("unable to get user's segments: %w", err)
	}

	// create map of user's segments
	userSegmentsMap := make(map[string]struct{}, len(userSegments))

	// add user's segments to the map
	for _, segment := range userSegments {
		userSegmentsMap[segment.Slug] = struct{}{}
	}

	var errorMessage strings.Builder
	var isError bool

	// check that user don't already have the segments we want to add
	for _, segment := range us.AddSegments {
		if _, ok := userSegmentsMap[segment.Slug]; ok {
			isError = true
			errorMessage.WriteString(fmt.Sprintf("user already have segment \"%v\"\n", segment.Slug))
		}
	}

	// check that user have the segments we want to remove
	for _, segment := range us.RemoveSegments {
		if _, ok := userSegmentsMap[segment.Slug]; !ok {
			isError = true
			errorMessage.WriteString(fmt.Sprintf("user don't have segment \"%v\"\n", segment.Slug))
		}
	}

	if isError {
		return fmt.Errorf("%w: %v", ErrIncorrectChangeUserSegmentsRequest, errorMessage.String())
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
		userSegmentsDB.RemoveSegments[i] = models.SegmentDeleteDB(segment)
	}

	// add the segments to the user
	err = s.db.ChangeUsersSegments(ctx, userSegmentsDB)
	if err != nil {
		return fmt.Errorf("unable to change user segments: %w", err)
	}

	return nil
}

func (s *SegmentifyDB) GetUsersSegments(ctx context.Context, userID int) (models.ActiveSegments, error) {
	segmentsDB, err := s.db.GetUsersSegments(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ActiveSegments{}, ErrNoUserData
		}
		return models.ActiveSegments{}, fmt.Errorf("unable to get user's segments: %w", err)
	}

	// in some cases GetUsersSegments returns empty slice instead of sql.ErrNoRows
	if len(segmentsDB) == 0 {
		return models.ActiveSegments{}, ErrNoUserData
	}

	segments := make(models.ActiveSegments, len(segmentsDB))
	for i, segmentDB := range segmentsDB {
		segments[i] = models.ActiveSegment{
			Slug: segmentDB.Slug,
		}
	}

	return segments, nil
}
