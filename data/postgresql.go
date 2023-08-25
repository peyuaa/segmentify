package data

import (
	"context"
	"fmt"
)

// SelectSegments returns a list of all segments from the database
func (s *SegmentifyDB) SelectSegments(ctx context.Context) (Segments, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, slug, is_deleted FROM segments")
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			s.l.Error("Unable to close rows", "error", err)
		}
	}()

	var segments Segments
	for rows.Next() {
		var segment Segment
		if err := rows.Scan(&segment.ID, &segment.Slug, &segment.IsDeleted); err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		segments = append(segments, segment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return segments, nil
}

// SelectSegmentByID returns a segment with given id from the database
func (s *SegmentifyDB) SelectSegmentByID(ctx context.Context, id int) (Segment, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Segment{}, fmt.Errorf("unable to begin transaction: %w", err)
	}

	var segment Segment
	err = tx.QueryRowContext(ctx, "SELECT id, slug, is_deleted FROM segments WHERE id = $1", id).
		Scan(&segment.ID, &segment.Slug, &segment.IsDeleted)
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			s.l.Error("Unable to rollback transaction", "error", rollErr)
		}
		return Segment{}, fmt.Errorf("unable to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return Segment{}, fmt.Errorf("unable to commit transaction: %w", err)
	}

	return segment, nil
}

// InsertSegment inserts segment with given slug into the database
func (s *SegmentifyDB) InsertSegment(ctx context.Context, slug string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO segments (slug) VALUES ($1)", slug)
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			s.l.Error("Unable to rollback transaction", "error", rollErr)
		}
		return fmt.Errorf("unable to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

// IsSegmentExists checks if segment with given slug exists in the database
// Returns true if segment exists, false otherwise
func (s *SegmentifyDB) IsSegmentExists(ctx context.Context, slug string) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM segments WHERE slug = $1", slug).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("unable to execute query: %w", err)
	}

	return count > 0, nil
}
