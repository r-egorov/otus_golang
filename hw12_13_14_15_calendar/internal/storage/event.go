package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	Title       string
	DateTime    time.Time
	Duration    time.Duration
	Description string
	OwnerID     uuid.UUID
	// NotifyBefore time.Duration
}

func (e *Event) String() string {
	return fmt.Sprintf("Event<ID: %s, Title: %s>",
		e.ID.String(), e.Title)
}
