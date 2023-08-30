package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/peyuaa/segmentify/db"
	"github.com/peyuaa/segmentify/models"

	"github.com/charmbracelet/log"
)

var (
	// ErrSegmentNotFound is an error returned when a segment can not be found in the database
	ErrSegmentNotFound = fmt.Errorf("segment not found")

	// ErrSegmentDeleted is an error returned when requested segment was marked as deleted
	ErrSegmentDeleted = fmt.Errorf("segment deleted")

	// ErrSegmentAlreadyExists is an error returned when a segment already exists in the database
	ErrSegmentAlreadyExists = fmt.Errorf("segment already exists")

	// ErrIncorrectChangeUserSegmentsRequest is an error returned when a request to change user segments is incorrect
	ErrIncorrectChangeUserSegmentsRequest = fmt.Errorf("incorrect change user segments request")

	// ErrNoUserData is an error returned when there is no user data about segments for given userID
	ErrNoUserData = fmt.Errorf("no user data about segments for given userID")

	// ErrNoUserHistoryData is an error returned when there is no user history data about segments for given userID
	// for specified period.
	ErrNoUserHistoryData = fmt.Errorf("no user history data about segments for given userID")
)

type SegmentifyDB struct {
	l  *log.Logger
	db *db.PostgresWrapper
}

func New(l *log.Logger, db *db.PostgresWrapper) *SegmentifyDB {
	return &SegmentifyDB{
		l:  l,
		db: db,
	}
}

func (s *SegmentifyDB) Add(ctx context.Context, segment models.CreateSegmentRequest) error {
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

func (s *SegmentifyDB) GetSegments(ctx context.Context) (models.Segments, error) {
	segmentsDB, err := s.db.SelectSegments(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get segments: %w", err)
	}

	segments := make(models.Segments, len(segmentsDB))
	for i, segmentDB := range segmentsDB {
		segments[i] = models.Segment(segmentDB)
	}

	return segments, nil
}

func (s *SegmentifyDB) GetSegmentBySlug(ctx context.Context, slug string) (models.Segment, error) {
	segmentDB, err := s.db.SelectSegmentBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Segment{}, ErrSegmentNotFound
		}
		return models.Segment{}, fmt.Errorf("unable to get segment by slug: %w", err)
	}

	segment := models.Segment(segmentDB)

	return segment, nil
}

func (s *SegmentifyDB) Delete(ctx context.Context, slug string) error {
	isDeleted, err := s.db.IsSegmentDeleted(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSegmentNotFound
		}
		return fmt.Errorf("unable to check whether segment is active: %w", err)
	}
	if isDeleted {
		return ErrSegmentNotFound
	}

	err = s.db.DeleteSegment(ctx, slug)
	if err != nil {
		return fmt.Errorf("unable to delete segment: %w", err)
	}
	return nil
}
