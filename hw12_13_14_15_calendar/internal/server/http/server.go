package internalhttp

import (
	"context"
	"errors"
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

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	return mux
}

func NewServer(logger server.Logger, app server.Application, host, port string) *Server {
	mux := newRouter()
	srv := &http.Server{Addr: host + ":" + port, Handler: loggingMiddleware(mux, logger)}
	return &Server{
		srv:  srv,
		app:  app,
		log:  logger,
		host: host,
		port: port,
	}
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatal("failed to start http server: " + err.Error())
		}
	}()
	s.log.Info(fmt.Sprintf("serving at %s", s.srv.Addr))
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("stopping server...")
	return s.srv.Shutdown(ctx)
}
