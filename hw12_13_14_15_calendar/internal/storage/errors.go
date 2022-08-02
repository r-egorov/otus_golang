package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ErrDateBusy struct {
	OwnerID uuid.UUID
	Date    time.Time
}

func (e ErrDateBusy) Error() string {
	return fmt.Sprintf(
		"user with ID <%s> has date: %s busy",
		e.OwnerID.String(), e.Date.Format("01-01-2006"),
	)
}

func NewErrDateBusy(ownerID uuid.UUID, date time.Time) *ErrDateBusy {
	return &ErrDateBusy{
		OwnerID: ownerID,
		Date:    date,
	}
}
