//go:build integration
// +build integration

package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

func TestSQLStorage_SaveEvent(t *testing.T) {
	s := New("postgres", "postgres", "calendar_test", "localhost", "5432")
	ctx := context.Background()
	require.NoError(t, s.Connect(ctx))

	t.Run("basic", func(t *testing.T) {
		defer truncateTable(t, s.db)
		loc, err := time.LoadLocation("Europe/Moscow")
		require.NoError(t, err)
		date := time.Date(2020, 1, 1, 13, 30, 0, 0, loc)
		event := storage.Event{
			ID:          uuid.New(),
			Title:       "Some title",
			DateTime:    date,
			Duration:    time.Hour * 2,
			Description: "Description",
			OwnerID:     uuid.New(),
		}
		gotEvent, err := s.SaveEvent(ctx, event)
		require.NoError(t, err)
		require.Equal(t, event, gotEvent)

		gotEvent, err = selectEvent(s.db, event.ID)
		require.NoError(t, err)
		require.Equal(t, event, gotEvent)
	})
}

func truncateTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`TRUNCATE TABLE events`)
	require.NoError(t, err)
}

func selectEvent(db *sql.DB, eventID uuid.UUID) (storage.Event, error) {
	row := db.QueryRow(`SELECT id, title, datetime, duration, description, owner_id FROM events WHERE id = $1`, eventID)

	var (
		event          storage.Event
		sqlDuration    sql.NullInt64
		sqlDescription sql.NullString
	)

	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.DateTime,
		&sqlDuration,
		&sqlDescription,
		&event.OwnerID,
	)

	switch err {
	case sql.ErrNoRows:
		return event, fmt.Errorf("no event with id %s", eventID.String())
	case nil:
		if sqlDuration.Valid {
			event.Duration = time.Duration(sqlDuration.Int64)
		}
		if sqlDescription.Valid {
			event.Description = sqlDescription.String
		}
		return event, nil
	default:
		return event, err
	}
}
