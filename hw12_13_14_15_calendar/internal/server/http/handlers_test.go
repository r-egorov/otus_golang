package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockApp struct {
	storage app.Storage
}

func newMockApp(s app.Storage) mockApp {
	return mockApp{s}
}

func (m *mockApp) SaveEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	event, err := m.storage.SaveEvent(ctx, event)
	if err != nil {
		return storage.Event{}, err
	}
	return event, nil
}

func (m *mockApp) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	event, err := m.storage.UpdateEvent(ctx, event)
	if err != nil {
		return storage.Event{}, err
	}
	return event, nil
}

func (m *mockApp) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	err := m.storage.DeleteEvent(ctx, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (m *mockApp) listEvents(
	ctx context.Context,
	listFunc func(context.Context, time.Time) ([]storage.Event, error),
	startPeriod time.Time,
) ([]storage.Event, error) {
	events, err := listFunc(ctx, startPeriod)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (m *mockApp) ListEventsDay(ctx context.Context, dayStart time.Time) ([]storage.Event, error) {
	events, err := m.listEvents(ctx, m.storage.ListEventsDay, dayStart)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (m *mockApp) ListEventsWeek(ctx context.Context, weekStart time.Time) ([]storage.Event, error) {
	events, err := m.listEvents(ctx, m.storage.ListEventsWeek, weekStart)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (m *mockApp) ListEventsMonth(ctx context.Context, monthStart time.Time) ([]storage.Event, error) {
	events, err := m.listEvents(ctx, m.storage.ListEventsMonth, monthStart)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func Test_CreateEvent(t *testing.T) {
	t.Run("creates event", func(t *testing.T) {
		s := memorystorage.New()
		a := newMockApp(s)
		mux := http.NewServeMux()
		mux.HandleFunc("/events", createEventsHandler(&a))

		datetime := time.Date(2022, time.Month(3), 1, 0, 0, 0, 0, time.UTC)
		expected := generateEvent(
			uuid.New(),
			"test created",
			datetime,
			time.Hour*2,
			"test description",
			uuid.New(),
		)

		reqBody := &bytes.Buffer{}
		err := json.NewEncoder(reqBody).Encode(expected)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/events", reqBody)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		require.Equal(t, http.StatusCreated, rr.Code)

		got := storage.Event{}
		err = json.NewDecoder(rr.Body).Decode(&got)
		require.NoError(t, err)

		expected.ID = got.ID
		require.Equal(t, expected, got)

		saved, err := s.ListEventsDay(context.Background(), datetime)
		require.NoError(t, err)
		require.Equal(t, 1, len(saved))
		require.Equal(t, expected, saved[0])
	})

	t.Run("it generates ID on the backend", func(t *testing.T) {
		s := memorystorage.New()
		a := newMockApp(s)
		mux := http.NewServeMux()
		mux.HandleFunc("/events", createEventsHandler(&a))

		reqBody := &bytes.Buffer{}
		reqBody.WriteString(`{
"title":"test created",
"datetime":"2022-03-01T00:00:00Z","duration":7200000000000,
"description":"test description",
"owner_id":"61b662df-7661-496f-8ada-8a04d1bfe78a"
}`)

		req, err := http.NewRequest("POST", "/events", reqBody)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		require.Equal(t, http.StatusCreated, rr.Code)

		got := storage.Event{}
		err = json.NewDecoder(rr.Body).Decode(&got)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, got.ID)
	})

	t.Run("method not allowed", func(t *testing.T) {
		s := memorystorage.New()
		a := newMockApp(s)
		mux := http.NewServeMux()
		mux.HandleFunc("/events", createEventsHandler(&a))

		req, err := http.NewRequest("GET", "/events", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})
}

func generateEvent(
	id uuid.UUID,
	title string,
	datetime time.Time,
	duration time.Duration,
	descritpion string,
	ownerID uuid.UUID,
) storage.Event {
	return storage.Event{
		ID:          id,
		Title:       title,
		DateTime:    datetime,
		Duration:    duration,
		Description: descritpion,
		OwnerID:     ownerID,
	}
}
