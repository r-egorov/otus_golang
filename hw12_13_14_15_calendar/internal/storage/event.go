package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	DateTime    time.Time     `json:"datetime"`
	Duration    time.Duration `json:"duration"`
	Description string        `json:"description"`
	OwnerID     uuid.UUID     `json:"owner_id"`
	// NotifyBefore time.Duration
}

func (e *Event) String() string {
	return fmt.Sprintf("Event<ID: %s, Title: %s>",
		e.ID.String(), e.Title)
}
