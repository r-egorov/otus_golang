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
		event := generateEvent()
		gotEvent, err := s.SaveEvent(ctx, event)
		require.NoError(t, err)
		require.Equal(t, event, gotEvent)

		gotEvent, err = selectEvent(s.db, event.ID)
		require.NoError(t, err)
		require.Equal(t, event, gotEvent)
	})

	t.Run("it returns err if ID is not unique", func(t *testing.T) {
		refEvent := generateEvent()
		_, err := s.SaveEvent(ctx, refEvent)
		require.NoError(t, err)

		sameIDEvent := generateEvent()
		sameIDEvent.ID = refEvent.ID
		_, err = s.SaveEvent(ctx, sameIDEvent)
		require.Error(t, err)
		var errIDNotUnique *storage.ErrIDNotUnique
		require.ErrorAs(t, err, &errIDNotUnique)

		eventInDB, err := selectEvent(s.db, refEvent.ID)
		require.NoError(t, err)
		require.Equal(t, refEvent, eventInDB)
	})

	t.Run("it returns err if date is busy", func(t *testing.T) {
		refEvent := generateEvent()
		_, err := s.SaveEvent(ctx, refEvent)
		require.NoError(t, err)

		sameDateEvent := generateEvent()
		sameDateEvent.DateTime = refEvent.DateTime
		sameDateEvent.OwnerID = refEvent.OwnerID
		_, err = s.SaveEvent(ctx, sameDateEvent)
		var errDateBusy *storage.ErrDateBusy
		require.ErrorAs(t, err, &errDateBusy)

		eventInDB, err := selectEvent(s.db, refEvent.ID)
		require.NoError(t, err)
		require.Equal(t, refEvent, eventInDB)
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

func generateEvent() storage.Event {
	return storage.Event{
		ID:          uuid.New(),
		Title:       "New Title",
		DateTime:    time.Date(2020, 1, 1, 13, 30, 0, 0, time.UTC),
		Duration:    2 * time.Hour,
		Description: "New description",
		OwnerID:     uuid.New(),
	}
}
