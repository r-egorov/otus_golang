package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	store map[uuid.UUID]storage.Event
	mu    sync.RWMutex
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) SaveEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roundedRefDatetime := event.DateTime.Round(time.Minute)
	for _, eventInStore := range s.store {
		roundedEventDatetime := eventInStore.DateTime.Round(time.Minute)
		isOwner := event.OwnerID == eventInStore.OwnerID
		dateBusy := roundedEventDatetime == roundedRefDatetime
		if isOwner && dateBusy {
			return storage.Event{}, storage.NewErrDateBusy(event.OwnerID, event.DateTime)
		}
	}

	s.store[event.ID] = event
	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[event.ID] = event
	return event, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, event.ID)
	return nil
}

func (s *Storage) ListEventsDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]storage.Event, 0)
	refYear, refMonth, refDay := day.Date()
	for _, event := range s.store {
		eventYear, eventMonth, eventDay := event.DateTime.Date()
		if eventYear == refYear && eventMonth == refMonth && eventDay == refDay {
			res = append(res, event)
		}
	}
	return res, nil
}

func (s *Storage) ListEventsWeek(ctx context.Context, weekStart time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]storage.Event, 0)
	refYear, refWeek := weekStart.ISOWeek()
	for _, event := range s.store {
		eventYear, eventWeek := event.DateTime.ISOWeek()
		if eventYear == refYear && eventWeek == refWeek {
			res = append(res, event)
		}
	}
	return res, nil
}

func (s *Storage) ListEventsMonth(ctx context.Context, monthStart time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]storage.Event, 0)
	refYear, refMonth, _ := monthStart.Date()
	for _, event := range s.store {
		eventYear, eventMonth, _ := event.DateTime.Date()
		if eventYear == refYear && eventMonth == refMonth {
			res = append(res, event)
		}
	}
	return res, nil
}

func New() *Storage {
	return &Storage{
		store: make(map[uuid.UUID]storage.Event),
	}
}
