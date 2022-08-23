package app

import (
	"context"
	"fmt"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
	"time"

	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
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

type App struct {
	Logg    Logger
	storage Storage
	conf    config.Config
}

func New(Logg Logger, conf config.Config) *App {
	var storage Storage
	switch conf.Storage.StorageType {
	case config.PSQLStorageType:
		storage = sqlstorage.New(
			conf.Storage.User,
			conf.Storage.Password,
			conf.Storage.DBName,
			conf.Storage.Host,
			conf.Storage.Port,
		)
	default:
		storage = memorystorage.New()
	}

	return &App{
		Logg:    Logg,
		storage: storage,
		conf:    conf,
	}
}

func (a *App) SaveEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	event, err := a.storage.SaveEvent(ctx, event)
	if err != nil {
		a.Logg.Error(
			fmt.Sprintf(
				"can't save %s, err: %v",
				event, err),
		)
		return storage.Event{}, err
	}
	a.Logg.Info(fmt.Sprintf("saved %s", event))
	return event, nil
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	event, err := a.storage.UpdateEvent(ctx, event)
	if err != nil {
		a.Logg.Error(
			fmt.Sprintf(
				"can't update %s, err: %v",
				event, err),
		)
		return storage.Event{}, err
	}
	a.Logg.Info(fmt.Sprintf("updated %s", event))
	return event, nil
}

func (a *App) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	err := a.storage.DeleteEvent(ctx, eventID)
	if err != nil {
		a.Logg.Error(
			fmt.Sprintf(
				"can't delete event, err: %v",
				err),
		)
		return err
	}
	a.Logg.Info(fmt.Sprintf("deleted event, ID: %s", eventID.String()))
	return nil
}

func (a *App) listEvents(
	ctx context.Context,
	listFunc func(context.Context, time.Time) ([]storage.Event, error),
	startPeriod time.Time,
) ([]storage.Event, error) {
	events, err := listFunc(ctx, startPeriod)
	if err != nil {
		a.Logg.Error(
			fmt.Sprintf(
				"can't get events list, err: %v",
				err),
		)
		return nil, err
	}
	return events, nil
}

func (a *App) ListEventsDay(ctx context.Context, dayStart time.Time) ([]storage.Event, error) {
	events, err := a.listEvents(ctx, a.storage.ListEventsDay, dayStart)
	if err != nil {
		return nil, err
	}
	a.Logg.Info(fmt.Sprintf("got events list for day: %s", dayStart.Format("2006-01-02 15:04:05")))
	return events, nil
}

func (a *App) ListEventsWeek(ctx context.Context, weekStart time.Time) ([]storage.Event, error) {
	events, err := a.listEvents(ctx, a.storage.ListEventsWeek, weekStart)
	if err != nil {
		return nil, err
	}
	a.Logg.Info(fmt.Sprintf("got events list for week: %s", weekStart.Format("2006-01-02 15:04:05")))
	return events, nil
}

func (a *App) ListEventsMonth(ctx context.Context, monthStart time.Time) ([]storage.Event, error) {
	events, err := a.listEvents(ctx, a.storage.ListEventsMonth, monthStart)
	if err != nil {
		return nil, err
	}
	a.Logg.Info(fmt.Sprintf("got events list for month: %s", monthStart.Format("2006-01-02 15:04:05")))
	return events, nil
}
