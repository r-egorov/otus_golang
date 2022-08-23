package internalhttp

import (
	"context"
	"fmt"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server"
	"net/http"
)

type Server struct {
	srv        *http.Server
	app        server.Application
	log        server.Logger
	host, port string
}

func NewServer(logger server.Logger, app server.Application, host, port string) *Server {
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
