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
		expected := generateEvent()
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

		refEvent := generateEvent()
		refEvent.DateTime = datetime
		refEvent.OwnerID = ownerID
		_, err := store.SaveEvent(ctx, refEvent)
		require.NoError(t, err)

		eventSameDate := generateEvent()
		eventSameDate.DateTime = datetime
		eventSameDate.OwnerID = ownerID
		_, err = store.SaveEvent(ctx, eventSameDate)
		require.Error(t, err)
		var errDateBusy *storage.ErrDateBusy
		require.ErrorAs(t, err, &errDateBusy)
		require.Equal(t, refEvent, store.store[refEvent.ID])
	})

	t.Run("it returns err if ID is not unique", func(t *testing.T) {
		ctx := context.Background()
		store := New()
		eventID := uuid.New()

		refEvent := generateEvent()
		refEvent.ID = eventID
		_, err := store.SaveEvent(ctx, refEvent)
		require.NoError(t, err)

		sameIDEvent := generateEvent()
		sameIDEvent.ID = eventID
		_, err = store.SaveEvent(ctx, sameIDEvent)
		require.Error(t, err)
		var errIDNotUnique *storage.ErrIDNotUnique
		require.ErrorAs(t, err, &errIDNotUnique)
		require.Equal(t, refEvent, store.store[refEvent.ID])
	})
}

func TestStorage_UpdateEvent(t *testing.T) {
	t.Run("it updates the event", func(t *testing.T) {
		ctx := context.Background()
		store := New()
		toUpdate := generateEvent()
		_, err := store.SaveEvent(ctx, toUpdate)
		require.NoError(t, err)
		require.Equal(t, toUpdate, store.store[toUpdate.ID])

		toUpdate.Title = "different event now"
		toUpdate.Duration = 15 * time.Minute
		toUpdate.DateTime = time.Date(2022, time.Month(3), 1, 12, 30, 0, 0, time.UTC)
		toUpdate.Description = "description is different"
		updated, err := store.UpdateEvent(ctx, toUpdate)
		require.NoError(t, err)
		require.Equal(t, toUpdate, updated)
		require.Equal(t, updated, store.store[updated.ID])
	})

	t.Run("it returns err if the date is busy", func(t *testing.T) {
		ctx := context.Background()
		store := New()
		datetime := time.Date(2022, time.Month(3), 1, 12, 30, 0, 0, time.UTC)
		ownerID := uuid.New()

		refEvent := generateEvent()
		refEvent.DateTime = datetime
		refEvent.OwnerID = ownerID
		_, err := store.SaveEvent(ctx, refEvent)
		require.NoError(t, err)

		storedEvent := generateEvent()
		storedEvent.OwnerID = ownerID
		_, err = store.SaveEvent(ctx, storedEvent)

		require.NoError(t, err)
		require.Equal(t, refEvent, store.store[refEvent.ID])
		require.Equal(t, storedEvent, store.store[storedEvent.ID])

		toUpdate := storedEvent
		toUpdate.DateTime = datetime
		_, err = store.UpdateEvent(ctx, toUpdate)
		require.Error(t, err)
		var errDateBusy *storage.ErrDateBusy
		require.ErrorAs(t, err, &errDateBusy)
		require.Equal(t, storedEvent, store.store[storedEvent.ID])
		require.NotEqual(t, toUpdate, store.store[storedEvent.ID])
	})

	t.Run("it returns err if ID is not found", func(t *testing.T) {
		ctx := context.Background()
		store := New()

		eventInStore := generateEvent()
		_, err := store.SaveEvent(ctx, eventInStore)
		require.NoError(t, err)

		newEvent := generateEvent()
		_, err = store.UpdateEvent(ctx, newEvent)
		require.Error(t, err)
		var errIDNotFound *storage.ErrIDNotFound
		require.ErrorAs(t, err, &errIDNotFound)
		_, exists := store.store[newEvent.ID]
		require.False(t, exists)
	})
}

func TestStorage_DeleteEvent(t *testing.T) {
	t.Run("it deletes the event", func(t *testing.T) {
		ctx := context.Background()
		store := New()

		eventInStore := generateEvent()
		_, err := store.SaveEvent(ctx, eventInStore)
		require.NoError(t, err)

		err = store.DeleteEvent(ctx, eventInStore.ID)
		require.NoError(t, err)
		_, exists := store.store[eventInStore.ID]
		require.False(t, exists)
	})

	t.Run("it returns err if ID is not found", func(t *testing.T) {
		ctx := context.Background()
		store := New()

		eventInStore := generateEvent()
		_, err := store.SaveEvent(ctx, eventInStore)
		require.NoError(t, err)

		err = store.DeleteEvent(ctx, uuid.New())
		require.Error(t, err)

		var errIDNotFound *storage.ErrIDNotFound
		require.ErrorAs(t, err, &errIDNotFound)

		_, exists := store.store[eventInStore.ID]
		require.True(t, exists)
	})
}

func TestStorage_ListEventsDay(t *testing.T) {
	t.Run("it returns list of events on a day", func(t *testing.T) {
		ctx := context.Background()
		store := New()
		date := time.Date(2022, time.Month(3), 1, 0, 0, 0, 0, time.UTC)
		morning := time.Date(2022, time.Month(3), 1, 6, 30, 0, 0, time.UTC)
		afternoon := time.Date(2022, time.Month(3), 1, 15, 30, 0, 0, time.UTC)
		evening := time.Date(2022, time.Month(3), 1, 21, 30, 0, 0, time.UTC)

		eventOne := generateEvent()
		eventOne.DateTime = morning
		_, err := store.SaveEvent(ctx, eventOne)
		require.NoError(t, err)

		eventTwo := generateEvent()
		eventTwo.DateTime = afternoon
		_, err = store.SaveEvent(ctx, eventTwo)
		require.NoError(t, err)

		eventThree := generateEvent()
		eventThree.DateTime = evening
		_, err = store.SaveEvent(ctx, eventThree)
		require.NoError(t, err)

		otherEvent := generateEvent()
		_, err = store.SaveEvent(ctx, otherEvent)
		require.NoError(t, err)

		expected := []storage.Event{eventOne, eventTwo, eventThree}
		got, err := store.ListEventsDay(ctx, date)
		require.NoError(t, err)
		assertEventListsEqual(t, expected, got)
	})

	t.Run("it returns empty list of events", func(t *testing.T) {
		ctx := context.Background()
		store := New()
		morning := time.Date(2022, time.Month(3), 1, 6, 30, 0, 0, time.UTC)
		afternoon := time.Date(2022, time.Month(3), 1, 15, 30, 0, 0, time.UTC)

		eventOne := generateEvent()
		eventOne.DateTime = morning
		_, err := store.SaveEvent(ctx, eventOne)
		require.NoError(t, err)

		eventTwo := generateEvent()
		eventTwo.DateTime = afternoon
		_, err = store.SaveEvent(ctx, eventTwo)
		require.NoError(t, err)

		date := time.Date(2022, time.Month(12), 12, 0, 0, 0, 0, time.UTC)
		got, err := store.ListEventsDay(ctx, date)
		require.NoError(t, err)
		require.Equal(t, []storage.Event{}, got)
	})
}

func generateEvent() storage.Event {
	return storage.Event{
		ID:          uuid.New(),
		Title:       "New Title",
		DateTime:    time.Now(),
		Duration:    2 * time.Hour,
		Description: "New description",
		OwnerID:     uuid.New(),
	}
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

	require.Equal(t, lmap, rmap)
}
