package storage

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	EventID    uuid.UUID `json:"eventId"`
	EventTitle string    `json:"eventTitle"`
	DateTime   time.Time `json:"datetime"`
	OwnerID    uuid.UUID `json:"ownerId"`
}

func (n *Notification) String() string {
	return fmt.Sprintf("Notification<EventID: %s, EventTitle: %s>",
		n.EventID.String(), n.EventTitle)
}
