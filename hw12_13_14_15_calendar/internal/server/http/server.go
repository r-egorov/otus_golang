package internalhttp

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"net/http"
	"time"
)

type Server struct {
	srv        *http.Server
	app        Application
	log        Logger
	host, port string
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

type Application interface {
	SaveEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	ListEventsDay(ctx context.Context, dayStart time.Time) ([]storage.Event, error)
	ListEventsWeek(ctx context.Context, weekStart time.Time) ([]storage.Event, error)
	ListEventsMonth(ctx context.Context, monthStart time.Time) ([]storage.Event, error)
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	srv := &http.Server{Addr: host + ":" + port, Handler: loggingMiddleware(mux, logger)}
	return &Server{
		srv:  srv,
		app:  app,
		log:  logger,
		host: host,
		port: port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	errChan := make(chan error)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	s.log.Info(fmt.Sprintf("serving at %s", s.srv.Addr))
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("stopping server...")
	return s.srv.Shutdown(ctx)
}
