package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/peyuaa/segmentify/models"

	"github.com/charmbracelet/log"
)

type PostgresWrapper struct {
	l  *log.Logger
	db *sql.DB
}

func New(l *log.Logger, db *sql.DB) *PostgresWrapper {
	return &PostgresWrapper{
		l:  l,
		db: db,
	}
}

// SelectSegments returns a list of all segments from the database
func (p *PostgresWrapper) SelectSegments(ctx context.Context) (models.SegmentsDB, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, slug, is_deleted FROM segments")
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			p.l.Error("Unable to close rows", "error", err)
		}
	}()

	var segments models.SegmentsDB
	for rows.Next() {
		var segment models.SegmentDB
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

func (p *PostgresWrapper) SelectSegmentBySlug(ctx context.Context, slug string) (models.SegmentDB, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return models.SegmentDB{}, fmt.Errorf("unable to begin transaction: %w", err)
	}

	var segment models.SegmentDB
	err = tx.QueryRowContext(ctx, "SELECT id, slug, is_deleted FROM segments WHERE slug = $1", slug).
		Scan(&segment.ID, &segment.Slug, &segment.IsDeleted)
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			p.l.Error("Unable to rollback transaction", "error", rollErr)
		}
		return models.SegmentDB{}, fmt.Errorf("unable to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return models.SegmentDB{}, fmt.Errorf("unable to commit transaction: %w", err)
	}

	return segment, nil
}

// InsertSegment inserts segment with given slug into the database
func (p *PostgresWrapper) InsertSegment(ctx context.Context, slug string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO segments (slug) VALUES ($1)", slug)
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			p.l.Error("Unable to rollback transaction", "error", rollErr)
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
func (p *PostgresWrapper) IsSegmentExists(ctx context.Context, slug string) (bool, error) {
	var count int
	err := p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM segments WHERE slug = $1", slug).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("unable to execute query: %w", err)
	}

	return count > 0, nil
}

func (p *PostgresWrapper) IsSegmentDeleted(ctx context.Context, slug string) (bool, error) {
	var isDeleted bool
	err := p.db.QueryRowContext(ctx, "SELECT is_deleted FROM segments WHERE slug = $1", slug).Scan(&isDeleted)
	if err != nil {
		return false, fmt.Errorf("unable to execute query: %w", err)
	}

	return isDeleted, nil
}

func (p *PostgresWrapper) DeleteSegment(ctx context.Context, slug string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE segments SET is_deleted = true WHERE slug = $1", slug)
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			p.l.Error("Unable to rollback transaction", "error", rollErr)
		}
		return fmt.Errorf("unable to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

// ChangeUsersSegments changes the segments of a user
// It calls addSegmentsToUser and deleteUserSegments and stores the segments addition and deletion history in one transaction
func (p *PostgresWrapper) ChangeUsersSegments(ctx context.Context, us models.UserSegmentsDB) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			rollErr := tx.Rollback()
			if rollErr != nil {
				p.l.Error("Unable to rollback transaction", "error", rollErr)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			err = fmt.Errorf("unable to commit transaction: %w", err)
		}
	}()

	// time of change
	t := time.Now()

	// add the segments to the user
	err = p.AddSegmentsToUser(ctx, tx, us.ID, us.AddSegments)
	if err != nil {
		return fmt.Errorf("unable to add segments to user: %w", err)
	}

	// add the segments to the user history
	err = p.AddSegmentInUsersHistory(ctx, tx, us.ID, us.AddSegments, t)
	if err != nil {
		return fmt.Errorf("unable to add segments to user history: %w", err)
	}

	// remove the segments from the user
	err = p.DeleteUserSegments(ctx, tx, us.ID, us.RemoveSegments)
	if err != nil {
		return fmt.Errorf("unable to delete segments from user: %w", err)
	}

	// add the deleted segments to the user history
	err = p.AddSegmentsRemoveDateInUserHistory(ctx, tx, us.ID, us.RemoveSegments, t)
	if err != nil {
		return fmt.Errorf("unable to add deleted segments to user history: %w", err)
	}

	return nil
}

// AddSegmentsToUser add segments to user using transaction tx
func (p *PostgresWrapper) AddSegmentsToUser(ctx context.Context, tx *sql.Tx, userID int, segments []models.SegmentAddDB) (err error) {
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO users_segments (user_id, slug, expiration_date) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %w", err)
	}
	defer func() {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			p.l.Error("Unable to close statement", "error", stmtErr)
		}
	}()

	for _, segment := range segments {
		_, err := stmt.ExecContext(ctx, userID, segment.Slug, segment.Expired)
		if err != nil {
			return fmt.Errorf("unable to execute query: %w", err)
		}
	}

	return nil
}

func (p *PostgresWrapper) AddSegmentInUsersHistory(ctx context.Context, tx *sql.Tx, userID int, segments []models.SegmentAddDB, time time.Time) error {
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO user_segment_history (user_id, segment_slug, date_added) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %w", err)
	}
	defer func() {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			p.l.Error("Unable to close statement", "error", stmtErr)
		}
	}()

	for _, segment := range segments {
		_, err := stmt.ExecContext(ctx, userID, segment.Slug, time)
		if err != nil {
			return fmt.Errorf("unable to execute query: %w", err)
		}
	}

	return nil
}

func (p *PostgresWrapper) AddSegmentsRemoveDateInUserHistory(ctx context.Context, tx *sql.Tx, userID int, segments []models.SegmentDeleteDB, time time.Time) error {
	stmt, err := tx.PrepareContext(ctx, "UPDATE user_segment_history SET date_removed = $1 WHERE user_id = $2 AND segment_slug = $3 AND date_removed IS NULL")
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %w", err)
	}
	defer func() {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			p.l.Error("Unable to close statement", "error", stmtErr)
		}
	}()

	for _, segment := range segments {
		_, err := stmt.ExecContext(ctx, time, userID, segment.Slug)
		if err != nil {
			return fmt.Errorf("unable to execute query: %w", err)
		}
	}

	return nil
}

// DeleteUserSegments deletes segments from user using transaction tx
func (p *PostgresWrapper) DeleteUserSegments(ctx context.Context, tx *sql.Tx, userID int, segments []models.SegmentDeleteDB) error {
	stmt, err := tx.PrepareContext(ctx, "DELETE FROM users_segments WHERE user_id = $1 AND slug = $2")
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %w", err)
	}
	defer func() {
		stmtErr := stmt.Close()
		if stmtErr != nil {
			p.l.Error("Unable to close statement", "error", stmtErr)
		}
	}()

	for _, segment := range segments {
		_, err := stmt.ExecContext(ctx, userID, segment.Slug)
		if err != nil {
			return fmt.Errorf("unable to execute query: %w", err)
		}
	}

	return nil
}
