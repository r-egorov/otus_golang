package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Detail string `json:"detail"`
}

type CreateEventRequest struct {
	Event storage.Event `json:"event"`
}

type CreateEventResponse struct {
	Event storage.Event `json:"event"`
}

func writeServerError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
}

func writeError(w http.ResponseWriter, err error, status int) {
	payload := ErrorResponse{err.Error()}

	body, err := json.Marshal(payload)
	if err != nil {
		writeServerError(w)
		return
	}

	w.WriteHeader(status)
	_, _ = fmt.Fprint(w, string(body))
}

func createEventsHandler(app server.Application) func(http.ResponseWriter, *http.Request) {
	postHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var event storage.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		event.ID = uuid.New()
		saved, err := app.SaveEvent(ctx, event)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(saved)
		if err != nil {
			writeServerError(w)
			return
		}
	}
	return postHandler
}
