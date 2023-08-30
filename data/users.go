package data

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/peyuaa/segmentify/models"
)

const (
	// history/userID/startDate/endDate
	historyDirTemplate = "history/%v/%v/%v"

	historyFileName = "history.csv"

	operationAdd    = "add"
	operationRemove = "remove"
)

func (s *SegmentifyDB) ChangeUserSegments(ctx context.Context, us models.UserSegmentsRequest) error {
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

func (s *SegmentifyDB) GetUserHistory(ctx context.Context, userID int, from, to time.Time) (filename string, err error) {
	history, err := s.db.GetUsersHistory(ctx, userID, from, to)

	switch {
	case err == nil:
		if len(history) == 0 {
			return filename, ErrNoUserHistoryData
		}
	case errors.Is(err, sql.ErrNoRows):
		return filename, ErrNoUserHistoryData
	default:
		return filename, fmt.Errorf("unable to get user's segments history: %w", err)
	}

	preparedHistory := s.prepareHistoryEntries(history, from, to)

	return s.writeCSV(preparedHistory, from, to)
}

func (s *SegmentifyDB) prepareHistoryEntries(db models.UserSegmentsHistoryDB, from, to time.Time) models.UserHistory {
	// len(db) is a minimum capacity of history, because every entry could be added and removed in the same period of time
	history := make(models.UserHistory, 0, len(db))

	for _, entry := range db {
		if entry.DateAdded.After(from) && entry.DateAdded.Before(to) {
			history = append(history, models.UserHistoryEntry{
				ID:        entry.ID,
				Slug:      entry.Slug,
				Operation: operationAdd,
				Date:      entry.DateAdded,
			})
		}
		if entry.DateRemoved.Valid && entry.DateRemoved.Time.After(from) && entry.DateRemoved.Time.Before(to) {
			history = append(history, models.UserHistoryEntry{
				ID:        entry.ID,
				Slug:      entry.Slug,
				Operation: operationRemove,
				Date:      entry.DateRemoved.Time,
			})
		}
	}

	// sort history by date ascending
	sort.Sort(history)

	return history
}

// writeCSV writes user's segments history to csv file
// and returns the path to the file and the error if any
func (s *SegmentifyDB) writeCSV(history models.UserHistory, from, to time.Time) (path string, err error) {
	// history[0] exists because we checked that len(history) > 0
	userID := history[0].ID

	// prepare records
	records := make([][]string, len(history))
	for i, segment := range history {
		records[i] = []string{
			strconv.Itoa(segment.ID),
			segment.Slug,
			segment.Operation,
			segment.Date.Format(time.RFC3339),
		}
	}

	dir := fmt.Sprintf(historyDirTemplate,
		userID, from.Format("2006-01-02"), to.Format("2006-01-02"))

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return path, fmt.Errorf("unable to create directory: %w", err)
	}

	// create csv file in directory "history/userID/startDate/endDate/history.csv"
	file, err := os.Create(dir + "/" + historyFileName)
	if err != nil {
		return path, fmt.Errorf("unable to create csv file: %w", err)
	}

	// write csv file
	err = csv.NewWriter(file).WriteAll(records)
	if err != nil {
		return path, fmt.Errorf("unable to write to csv file: %w", err)
	}

	return file.Name(), nil
}
