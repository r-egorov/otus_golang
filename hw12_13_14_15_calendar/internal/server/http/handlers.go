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

type UpdateEventRequest struct {
	Event storage.Event `json:"event"`
}

type UpdateEventResponse struct {
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

func eventsHandler(app server.Application) func(http.ResponseWriter, *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			createEventsHandler(app, w, r)
		case "PATCH":
			updateEventsHandler(app, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
	return handler
}

func createEventsHandler(app server.Application, w http.ResponseWriter, r *http.Request) {
	var request CreateEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	event := request.Event
	event.ID = uuid.New()
	saved, err := app.SaveEvent(ctx, event)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	response := CreateEventResponse{Event: saved}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeServerError(w)
		return
	}
}

func updateEventsHandler(app server.Application, w http.ResponseWriter, r *http.Request) {
	var request UpdateEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	event := request.Event
	updated, err := app.UpdateEvent(ctx, event)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	response := UpdateEventResponse{Event: updated}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeServerError(w)
		return
	}
}
