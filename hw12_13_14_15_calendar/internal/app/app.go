package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	SaveEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	ListEventsDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	ListEventsWeek(ctx context.Context, weekStart time.Time) ([]storage.Event, error)
	ListEventsMonth(ctx context.Context, monthStart time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
