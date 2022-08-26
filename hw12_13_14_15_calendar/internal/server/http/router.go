package internalhttp

import (
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server"
	"net/http"
)

type MyRouter struct {
	mux *http.ServeMux
	app server.Application
}

func newRouter(app server.Application) MyRouter {
	mux := http.NewServeMux()

	mux.HandleFunc("/events", createEventsHandler(app))

	router := MyRouter{
		mux: mux,
		app: app,
	}
	return router
}

func (m MyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}
