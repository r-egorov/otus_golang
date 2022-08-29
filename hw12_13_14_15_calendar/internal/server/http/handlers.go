package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"net/http"
	"time"
)

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
			createEventHandler(app, w, r)
		case "PATCH":
			updateEventHandler(app, w, r)
		case "DELETE":
			deleteEventHandler(app, w, r)
		case "GET":
			getEventsHandler(app, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
	return handler
}

func createEventHandler(app server.Application, w http.ResponseWriter, r *http.Request) {
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

func updateEventHandler(app server.Application, w http.ResponseWriter, r *http.Request) {
	var request UpdateEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	updated, err := app.UpdateEvent(ctx, request.Event)
	if err != nil {
		status := http.StatusBadRequest

		var errIDNotFound *storage.ErrIDNotFound
		if errors.As(err, &errIDNotFound) {
			status = http.StatusNotFound
		}

		writeError(w, err, status)
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

func deleteEventHandler(app server.Application, w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = app.DeleteEvent(ctx, id)
	if err != nil {
		status := http.StatusBadRequest

		var errIDNotFound *storage.ErrIDNotFound
		if errors.As(err, &errIDNotFound) {
			status = http.StatusNotFound
		}

		writeError(w, err, status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getEventsHandler(app server.Application, w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startPeriodStr := r.URL.Query().Get("datetime")
	if startPeriodStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startPeriod, err := time.Parse(time.RFC3339Nano, startPeriodStr)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	var listFunc func(context.Context, time.Time) ([]storage.Event, error)

	switch period {
	case "day":
		listFunc = app.ListEventsDay
	case "week":
		listFunc = app.ListEventsWeek
	case "month":
		listFunc = app.ListEventsMonth
	default:
		writeError(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	events, err := listFunc(ctx, startPeriod)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
	}

	response := GetEventsResponse{Events: events}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeServerError(w)
		return
	}
}
