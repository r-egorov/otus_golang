package storage

import (
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
