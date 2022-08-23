package app

import (
	"context"
	"fmt"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
	internalhttp "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
	"os"
	"os/signal"
	"syscall"
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
	ctx     context.Context
	logg    Logger
	storage Storage
	conf    config.Config
}

func New(logg Logger, conf config.Config) *App {
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
		ctx:     context.Background(),
		logg:    logg,
		storage: storage,
		conf:    conf,
	}
}

func (a *App) Run() {
	ctx, cancel := signal.NotifyContext(a.ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	httpserver := internalhttp.NewServer(a.logg, a, a.conf.Server.Host, a.conf.Server.Port)

	serverStopped := make(chan struct{})
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpserver.Stop(ctx); err != nil {
			a.logg.Error("failed to stop http server: " + err.Error())
		}
		serverStopped <- struct{}{}
	}()

	a.logg.Info("calendar is running...")

	if err := httpserver.Start(ctx); err != nil {
		a.logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
	<-serverStopped
}

func (a *App) CreateEvent(ctx context.Context, id uuid.UUID, title string) (storage.Event, error) {
	event, err := a.storage.SaveEvent(ctx, storage.Event{ID: id, Title: title})
	if err != nil {
		a.logg.Error(
			fmt.Sprintf(
				"can't create event ID: %s, Title: %s, err: %v",
				id, title, err),
		)
		return storage.Event{}, err
	}
	a.logg.Info(fmt.Sprintf("created event ID: %s, Title: %s", id, title))
	return event, nil
}

// TODO
