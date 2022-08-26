package internalhttp

import "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"

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
