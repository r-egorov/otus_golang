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

type ErrIDNotUnique struct {
	EventID uuid.UUID
}

func (e ErrIDNotUnique) Error() string {
	return fmt.Sprintf(
		"event with ID <%s> already exists",
		e.EventID.String(),
	)
}

func NewErrIDNotUnique(eventID uuid.UUID) *ErrIDNotUnique {
	return &ErrIDNotUnique{
		EventID: eventID,
	}
}

type ErrIDNotFound struct {
	EventID uuid.UUID
}

func (e ErrIDNotFound) Error() string {
	return fmt.Sprintf(
		"event with ID <%s> not found",
		e.EventID.String(),
	)
}

func NewErrIDNotFound(eventID uuid.UUID) *ErrIDNotFound {
	return &ErrIDNotFound{
		EventID: eventID,
	}
}
