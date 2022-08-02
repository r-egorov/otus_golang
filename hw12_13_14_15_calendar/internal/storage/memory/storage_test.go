package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage_SaveEvent(t *testing.T) {
	t.Run("it saves the event", func(t *testing.T) {
		ctx := context.Background()
		store := New()
		expected := storage.Event{
			ID:          uuid.New(),
			Title:       "Some Title",
			DateTime:    time.Date(2022, time.Month(1), 1, 12, 30, 0, 0, time.UTC),
			Duration:    2 * time.Hour,
			Description: "Event description",
			OwnerID:     uuid.New(),
		}
		got, err := store.SaveEvent(ctx, expected)
		require.NoError(t, err)
		require.Equal(t, expected, got)
		require.Equal(t, got, store.store[got.ID])
	})

	t.Run("it returns err if the date is busy", func(t *testing.T) {
		ctx := context.Background()
		store := New()
		ownerID := uuid.New()
		datetime := time.Date(2022, time.Month(1), 1, 12, 30, 0, 0, time.UTC)

		event := storage.Event{
			ID:          uuid.New(),
			Title:       "Some Title",
			DateTime:    datetime,
			Duration:    2 * time.Hour,
			Description: "Event description",
			OwnerID:     ownerID,
		}
		_, err := store.SaveEvent(ctx, event)
		require.NoError(t, err)

		eventSameDate := storage.Event{
			ID:          uuid.New(),
			Title:       "New Title",
			DateTime:    datetime,
			Duration:    4 * time.Hour,
			Description: "New description",
			OwnerID:     ownerID,
		}

		_, err = store.SaveEvent(ctx, eventSameDate)
		require.Error(t, err)
		var errDateBusy *storage.ErrDateBusy
		require.ErrorAs(t, err, &errDateBusy)
	})
}
