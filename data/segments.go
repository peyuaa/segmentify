package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
)

var (
	// ErrSegmentNotFound is an error raised when a segment can not be found in the database
	ErrSegmentNotFound = fmt.Errorf("segment not found")

	// ErrSegmentAlreadyExists is an error raised when a segment already exists in the database
	ErrSegmentAlreadyExists = fmt.Errorf("segment already exists")
)

// Segment defines the structure for an API segment
type Segment struct {
	// the id for the segment
	ID int `json:"id"` // Unique identifier for the segment

	// the segment's slug
	Slug string `json:"slug" validate:"required,min=5,max=50"`

	// is the segment deleted
	IsDeleted bool `json:"is_deleted"`
}

// Segments defines a slice of Segment
type Segments []Segment

type SegmentifyDB struct {
	l  *log.Logger
	db *sql.DB
}

func New(l *log.Logger, db *sql.DB) *SegmentifyDB {
	return &SegmentifyDB{
		l:  l,
		db: db,
	}
}

func (s *SegmentifyDB) Add(ctx context.Context, segment Segment) error {
	exists, err := s.isSegmentExists(ctx, segment.Slug)
	if err != nil {
		return fmt.Errorf("unable to check segment existence: %w", err)
	}
	if exists {
		return ErrSegmentAlreadyExists
	}

	err = s.insertSegment(ctx, segment.Slug)
	if err != nil {
		return fmt.Errorf("unable to insert segment: %w", err)
	}

	return nil
}

func (s *SegmentifyDB) GetSegments(ctx context.Context) (Segments, error) {
	segments, err := s.selectSegments(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get segments: %w", err)
	}

	return segments, nil
}

func (s *SegmentifyDB) GetSegmentBySlug(ctx context.Context, slug string) (Segment, error) {
	segment, err := s.selectSegmentBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return segment, ErrSegmentNotFound
		}
		return segment, fmt.Errorf("unable to get segment by slug: %w", err)
	}

	return segment, nil
}

func (s *SegmentifyDB) Delete(ctx context.Context, slug string) error {
	isDeleted, err := s.isSegmentDeleted(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSegmentNotFound
		}
		return fmt.Errorf("unable to check whether segment is active: %w", err)
	}
	if isDeleted {
		return ErrSegmentNotFound
	}

	err = s.deleteSegment(ctx, slug)
	if err != nil {
		return fmt.Errorf("unable to delete segment: %w", err)
	}
	return nil
}
