package internalhttp

import (
	"net/http"

	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server"
)

type MyRouter struct {
	mux *http.ServeMux
	app server.Application
}

func newRouter(app server.Application) MyRouter {
	mux := http.NewServeMux()

	mux.HandleFunc("/events", eventsHandler(app))

	router := MyRouter{
		mux: mux,
		app: app,
	}
	return router
}

func (m MyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}
