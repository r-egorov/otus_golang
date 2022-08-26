package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Detail string `json:"detail"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func createEventsHandler(app server.Application, method string) func(http.ResponseWriter, *http.Request) {
	postHandler := func(w http.ResponseWriter, r *http.Request) {
		var event storage.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			payload := ErrorResponse{err.Error()}

			body, err := json.Marshal(payload)
			if err != nil {
				status := http.StatusInternalServerError
				http.Error(w, http.StatusText(status), status)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, string(body))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		saved, err := app.SaveEvent(ctx, event)
		if err != nil {
			payload := ErrorResponse{err.Error()}

			body, err := json.Marshal(payload)
			if err != nil {
				status := http.StatusInternalServerError
				http.Error(w, http.StatusText(status), status)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, string(body))
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(saved)
		if err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	}

	switch method {
	case "POST":
		return postHandler
	}

	return hello
}
