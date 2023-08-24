package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/charmbracelet/log"
)

type Segmentify struct {
	l  *log.Logger
	db *sql.DB
}

func New(l *log.Logger, db *sql.DB) *Segmentify {
	return &Segmentify{
		l:  l,
		db: db,
	}
}

func (s *Segmentify) InsertSegment(ctx context.Context, slug string) error {
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

func (s *Segmentify) IsSegmentExists(ctx context.Context, slug string) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM segments WHERE slug = $1", slug).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("unable to execute query: %w", err)
	}

	return count > 0, nil
}
