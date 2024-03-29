//go:build integration
// +build integration

package sqlstorage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

const pgHost = "172.22.208.1"

func TestSQLStorage_SaveEvent(t *testing.T) {
	s := New("postgres", "postgres", "calendar_test", pgHost, "5433")
	ctx := context.Background()
	require.NoError(t, s.Connect(ctx))
	defer func() {
		require.NoError(t, s.Close(ctx))
	}()

	t.Run("basic", func(t *testing.T) {
		defer truncateTable(t, s.db)
		event := generateEvent()
		gotEvent, err := s.SaveEvent(ctx, event)
		event.ID = gotEvent.ID
		require.NoError(t, err)
		require.Equal(t, event, gotEvent)

		assertEventInDB(t, s.db, event)
	})

	t.Run("it returns err if ID is not unique", func(t *testing.T) {
		defer truncateTable(t, s.db)
		refEvent := generateEvent()
		refEvent, err := s.SaveEvent(ctx, refEvent)
		require.NoError(t, err)

		sameIDEvent := generateEvent()
		sameIDEvent.ID = refEvent.ID
		_, err = s.SaveEvent(ctx, sameIDEvent)
		require.Error(t, err)
		var errIDNotUnique *storage.ErrIDNotUnique
		require.ErrorAs(t, err, &errIDNotUnique)

		assertEventInDB(t, s.db, refEvent)
	})

	t.Run("it returns err if date is busy", func(t *testing.T) {
		defer truncateTable(t, s.db)
		refEvent := generateEvent()
		refEvent, err := s.SaveEvent(ctx, refEvent)
		require.NoError(t, err)

		sameDateEvent := generateEvent()
		sameDateEvent.DateTime = refEvent.DateTime
		sameDateEvent.OwnerID = refEvent.OwnerID
		_, err = s.SaveEvent(ctx, sameDateEvent)
		var errDateBusy *storage.ErrDateBusy
		require.ErrorAs(t, err, &errDateBusy)

		assertEventInDB(t, s.db, refEvent)
	})
}

func TestStorage_UpdateEvent(t *testing.T) {
	s := New("postgres", "postgres", "calendar_test", pgHost, "5433")
	ctx := context.Background()
	require.NoError(t, s.Connect(ctx))
	defer func() {
		require.NoError(t, s.Close(ctx))
	}()

	t.Run("basic", func(t *testing.T) {
		defer truncateTable(t, s.db)
		event := generateEvent()
		_, err := s.SaveEvent(ctx, event)
		require.NoError(t, err)

		event.Title = "WOW THAT'S NEW TITLE"
		event.Duration = 16 * time.Minute
		event.Description = "DESCRIPTION? ALSO NEW!"

		updEvent, err := s.UpdateEvent(ctx, event)
		require.NoError(t, err)
		require.Equal(t, event, updEvent)

		assertEventInDB(t, s.db, event)
	})

	t.Run("it returns err if the date is busy", func(t *testing.T) {
		defer truncateTable(t, s.db)
		busyEvent := generateEvent()
		_, err := s.SaveEvent(ctx, busyEvent)
		require.NoError(t, err)

		newEvent := generateEvent()
		newEvent.DateTime = newEvent.DateTime.Add(time.Hour * 48)
		newEvent.OwnerID = busyEvent.OwnerID
		_, err = s.SaveEvent(ctx, newEvent)
		require.NoError(t, err)

		newEventUpdated := newEvent
		newEventUpdated.DateTime = busyEvent.DateTime
		_, err = s.UpdateEvent(ctx, newEventUpdated)
		var errDateBusy *storage.ErrDateBusy
		require.ErrorAs(t, err, &errDateBusy)

		assertEventInDB(t, s.db, newEvent)
	})

	t.Run("it returns err if ID is not found", func(t *testing.T) {
		defer truncateTable(t, s.db)
		newEvent := generateEvent()
		_, err := s.UpdateEvent(ctx, newEvent)
		require.Error(t, err)
		var errIDNotFound *storage.ErrIDNotFound
		require.ErrorAs(t, err, &errIDNotFound)
	})
}

func TestStorage_DeleteEvent(t *testing.T) {
	s := New("postgres", "postgres", "calendar_test", pgHost, "5433")
	ctx := context.Background()
	require.NoError(t, s.Connect(ctx))
	defer func() {
		require.NoError(t, s.Close(ctx))
	}()

	t.Run("basic", func(t *testing.T) {
		defer truncateTable(t, s.db)
		eventInDB := generateEvent()
		_, err := s.SaveEvent(ctx, eventInDB)
		require.NoError(t, err)

		err = s.DeleteEvent(ctx, eventInDB.ID)
		require.NoError(t, err)

		_, err = selectEvent(s.db, eventInDB.ID)
		require.Error(t, err)
		var errIDNotFound *storage.ErrIDNotFound
		require.ErrorAs(t, err, &errIDNotFound)
	})

	t.Run("it returns err if id not found", func(t *testing.T) {
		defer truncateTable(t, s.db)
		err := s.DeleteEvent(ctx, uuid.New())
		require.Error(t, err)
		var errIDNotFound *storage.ErrIDNotFound
		require.ErrorAs(t, err, &errIDNotFound)
	})
}

func TestStorage_ListEventsDay(t *testing.T) {
	s := New("postgres", "postgres", "calendar_test", pgHost, "5433")
	ctx := context.Background()
	require.NoError(t, s.Connect(ctx))
	defer func() {
		require.NoError(t, s.Close(ctx))
	}()

	t.Run("it returns list of events on a day", func(t *testing.T) {
		defer truncateTable(t, s.db)

		date := time.Date(2022, time.Month(3), 16, 0, 0, 0, 0, time.UTC)
		morning := time.Date(2022, time.Month(3), 16, 6, 30, 0, 0, time.UTC)
		afternoon := time.Date(2022, time.Month(3), 16, 15, 30, 0, 0, time.UTC)
		evening := time.Date(2022, time.Month(3), 16, 21, 30, 0, 0, time.UTC)

		eventOne := generateEvent()
		eventOne.DateTime = morning
		_, err := s.SaveEvent(ctx, eventOne)
		require.NoError(t, err)

		eventTwo := generateEvent()
		eventTwo.DateTime = afternoon
		_, err = s.SaveEvent(ctx, eventTwo)
		require.NoError(t, err)

		eventThree := generateEvent()
		eventThree.DateTime = evening
		_, err = s.SaveEvent(ctx, eventThree)
		require.NoError(t, err)

		otherEvent := generateEvent()
		_, err = s.SaveEvent(ctx, otherEvent)
		require.NoError(t, err)

		expected := []storage.Event{eventOne, eventTwo, eventThree}
		got, err := s.ListEventsDay(ctx, date)
		require.NoError(t, err)
		assertEventListsEqual(t, expected, got)
	})

	t.Run("it returns empty list of events", func(t *testing.T) {
		defer truncateTable(t, s.db)

		morning := time.Date(2022, time.Month(3), 1, 6, 30, 0, 0, time.UTC)
		afternoon := time.Date(2022, time.Month(3), 1, 15, 30, 0, 0, time.UTC)

		eventOne := generateEvent()
		eventOne.DateTime = morning
		_, err := s.SaveEvent(ctx, eventOne)
		require.NoError(t, err)

		eventTwo := generateEvent()
		eventTwo.DateTime = afternoon
		_, err = s.SaveEvent(ctx, eventTwo)
		require.NoError(t, err)

		date := time.Date(2022, time.Month(12), 12, 0, 0, 0, 0, time.UTC)
		got, err := s.ListEventsDay(ctx, date)
		require.NoError(t, err)
		require.Equal(t, []storage.Event{}, got)
	})
}

func TestStorage_ListEventsWeek(t *testing.T) {
	s := New("postgres", "postgres", "calendar_test", pgHost, "5433")
	ctx := context.Background()
	require.NoError(t, s.Connect(ctx))
	defer func() {
		require.NoError(t, s.Close(ctx))
	}()

	t.Run("it returns list of events in a week", func(t *testing.T) {
		defer truncateTable(t, s.db)

		weekStart := time.Date(2022, time.Month(3), 7, 0, 0, 0, 0, time.UTC)
		thursday := time.Date(2022, time.Month(3), 10, 6, 30, 0, 0, time.UTC)
		friday := time.Date(2022, time.Month(3), 11, 15, 30, 0, 0, time.UTC)
		tuesday := time.Date(2022, time.Month(3), 8, 21, 30, 0, 0, time.UTC)

		eventOne := generateEvent()
		eventOne.DateTime = tuesday
		_, err := s.SaveEvent(ctx, eventOne)
		require.NoError(t, err)

		eventTwo := generateEvent()
		eventTwo.DateTime = friday
		_, err = s.SaveEvent(ctx, eventTwo)
		require.NoError(t, err)

		eventThree := generateEvent()
		eventThree.DateTime = thursday
		_, err = s.SaveEvent(ctx, eventThree)
		require.NoError(t, err)

		otherEvent := generateEvent()
		_, err = s.SaveEvent(ctx, otherEvent)
		require.NoError(t, err)

		expected := []storage.Event{eventOne, eventTwo, eventThree}
		got, err := s.ListEventsWeek(ctx, weekStart)
		require.NoError(t, err)
		assertEventListsEqual(t, expected, got)
	})

	t.Run("it returns empty list of events", func(t *testing.T) {
		defer truncateTable(t, s.db)

		thursday := time.Date(2022, time.Month(3), 10, 6, 30, 0, 0, time.UTC)
		friday := time.Date(2022, time.Month(3), 11, 15, 30, 0, 0, time.UTC)

		eventOne := generateEvent()
		eventOne.DateTime = thursday
		_, err := s.SaveEvent(ctx, eventOne)
		require.NoError(t, err)

		eventTwo := generateEvent()
		eventTwo.DateTime = friday
		_, err = s.SaveEvent(ctx, eventTwo)
		require.NoError(t, err)

		weekStart := time.Date(2022, time.Month(10), 10, 0, 0, 0, 0, time.UTC)
		got, err := s.ListEventsDay(ctx, weekStart)
		require.NoError(t, err)
		require.Equal(t, []storage.Event{}, got)
	})
}

func TestStorage_ListEventsMonth(t *testing.T) {
	s := New("postgres", "postgres", "calendar_test", pgHost, "5433")
	ctx := context.Background()
	require.NoError(t, s.Connect(ctx))
	defer func() {
		require.NoError(t, s.Close(ctx))
	}()

	t.Run("it returns list of events in a month", func(t *testing.T) {
		defer truncateTable(t, s.db)

		date := time.Date(2022, time.Month(3), 1, 0, 0, 0, 0, time.UTC)
		tenth := time.Date(2022, time.Month(3), 10, 6, 30, 0, 0, time.UTC)
		fifteenth := time.Date(2022, time.Month(3), 15, 15, 30, 0, 0, time.UTC)
		twentyfifth := time.Date(2022, time.Month(3), 25, 21, 30, 0, 0, time.UTC)

		eventOne := generateEvent()
		eventOne.DateTime = tenth
		_, err := s.SaveEvent(ctx, eventOne)
		require.NoError(t, err)

		eventTwo := generateEvent()
		eventTwo.DateTime = fifteenth
		_, err = s.SaveEvent(ctx, eventTwo)
		require.NoError(t, err)

		eventThree := generateEvent()
		eventThree.DateTime = twentyfifth
		_, err = s.SaveEvent(ctx, eventThree)
		require.NoError(t, err)

		otherEvent := generateEvent()
		_, err = s.SaveEvent(ctx, otherEvent)
		require.NoError(t, err)

		expected := []storage.Event{eventOne, eventTwo, eventThree}
		got, err := s.ListEventsMonth(ctx, date)
		require.NoError(t, err)
		assertEventListsEqual(t, expected, got)
	})

	t.Run("it returns empty list of events", func(t *testing.T) {
		defer truncateTable(t, s.db)

		tenth := time.Date(2022, time.Month(3), 10, 6, 30, 0, 0, time.UTC)
		fifteenth := time.Date(2022, time.Month(3), 15, 15, 30, 0, 0, time.UTC)

		eventOne := generateEvent()
		eventOne.DateTime = tenth
		_, err := s.SaveEvent(ctx, eventOne)
		require.NoError(t, err)

		eventTwo := generateEvent()
		eventTwo.DateTime = fifteenth
		_, err = s.SaveEvent(ctx, eventTwo)
		require.NoError(t, err)

		date := time.Date(2022, time.Month(12), 1, 0, 0, 0, 0, time.UTC)
		got, err := s.ListEventsDay(ctx, date)
		require.NoError(t, err)
		require.Equal(t, []storage.Event{}, got)
	})
}

func TestStorage_ListEventsToNotify(t *testing.T) {
	s := New("postgres", "postgres", "calendar_test", pgHost, "5433")
	ctx := context.Background()
	require.NoError(t, s.Connect(ctx))
	defer func() {
		require.NoError(t, s.Close(ctx))
	}()

	t.Run("it returns list of events in a month", func(t *testing.T) {
		defer truncateTable(t, s.db)

		datesMap := make(map[string]time.Time)
		datesMap["now"] = time.Date(2022, time.Month(1), 10, 12, 30, 0, 0, time.UTC)
		datesMap["inFiveMinutes"] = time.Date(2022, time.Month(1), 10, 12, 35, 0, 0, time.UTC)
		datesMap["inTenMinutes"] = time.Date(2022, time.Month(1), 10, 12, 40, 0, 0, time.UTC)
		datesMap["inFifteenMinutes"] = time.Date(2022, time.Month(1), 10, 12, 45, 0, 0, time.UTC)
		datesMap["inHour"] = time.Date(2022, time.Month(1), 10, 13, 30, 0, 0, time.UTC)

		events := []storage.Event{}
		for k, v := range datesMap {
			event := generateEvent()
			event.Title = k
			event.DateTime = v
			savedEvent, err := s.SaveEvent(ctx, event)
			require.NoError(t, err)
			events = append(events, savedEvent)
		}

		notifications, err := s.ListToNotify(ctx, datesMap["now"], 10*time.Minute, 5*time.Minute)
		require.NoError(t, err)
		require.Equal(t, 3, len(notifications))

		expected := make(map[string]struct{})
		expected["inFiveMinutes"] = struct{}{}
		expected["inTenMinutes"] = struct{}{}
		expected["inFifteenMinutes"] = struct{}{}

		got := make(map[string]struct{})
		for _, n := range notifications {
			got[n.EventTitle] = struct{}{}
		}

		require.Equal(t, expected, got)
	})
}

func TestStorage_ClearOlderThanYear(t *testing.T) {
	s := New("postgres", "postgres", "calendar_test", pgHost, "5433")
	ctx := context.Background()
	require.NoError(t, s.Connect(ctx))
	defer func() {
		require.NoError(t, s.Close(ctx))
	}()

	t.Run("basic", func(t *testing.T) {
		defer truncateTable(t, s.db)
		oldDate := time.Date(2020, 1, 1, 13, 30, 0, 0, time.UTC)
		eventInDB := generateEvent()
		eventInDB.DateTime = oldDate
		_, err := s.SaveEvent(ctx, eventInDB)
		require.NoError(t, err)

		freshEvent := generateEvent()
		freshEvent.DateTime = time.Now().UTC()
		_, err = s.SaveEvent(ctx, freshEvent)
		require.NoError(t, err)

		err = s.ClearOlderThanYear(ctx)
		require.NoError(t, err)

		_, err = selectEvent(s.db, eventInDB.ID)
		require.Error(t, err)
		var errIDNotFound *storage.ErrIDNotFound
		require.ErrorAs(t, err, &errIDNotFound)

		_, err = selectEvent(s.db, freshEvent.ID)
		require.NoError(t, err)
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
		return event, storage.NewErrIDNotFound(eventID)
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

func assertEventInDB(t *testing.T, db *sql.DB, refEvent storage.Event) {
	eventInDB, err := selectEvent(db, refEvent.ID)
	require.NoError(t, err)
	require.Equal(t, refEvent, eventInDB)
}

func assertEventListsEqual(t *testing.T, l, r []storage.Event) {
	t.Helper()

	mapFromSlice := func(sl []storage.Event) map[uuid.UUID]storage.Event {
		res := make(map[uuid.UUID]storage.Event)
		for _, event := range sl {
			res[event.ID] = event
		}
		return res
	}

	lmap := mapFromSlice(l)
	rmap := mapFromSlice(r)

	require.Equal(t, len(lmap), len(rmap))
	for k, lv := range lmap {
		rv := rmap[k]
		require.Equal(t, lv.ID, rv.ID)
		require.Equal(t, lv.Title, rv.Title)
		require.Equal(t, lv.Duration, rv.Duration)
		require.True(t, lv.DateTime.Equal(rv.DateTime))
		require.Equal(t, lv.Description, rv.Description)
		require.Equal(t, lv.OwnerID, rv.OwnerID)
	}
}

func assertNotificationListsEqual(t *testing.T, l, r []storage.Notification) {
	t.Helper()

	mapFromSlice := func(sl []storage.Notification) map[uuid.UUID]storage.Notification {
		res := make(map[uuid.UUID]storage.Notification)
		for _, notification := range sl {
			res[notification.EventID] = notification
		}
		return res
	}

	lmap := mapFromSlice(l)
	rmap := mapFromSlice(r)

	require.Equal(t, lmap, rmap)
}
