package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type Storage struct {
	User, Password, DBName, Host, Port string
	db                                 *sql.DB
}

func New(user, password, dbName, host, port string) *Storage {
	return &Storage{
		User:     user,
		Password: password,
		DBName:   dbName,
		Host:     host,
		Port:     port,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		s.User, s.Password, s.Host, s.Port, s.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) SaveEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	fmt.Println("here")
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return event, err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO events (id, title, datetime, duration, description, owner_id) VALUES ($1, $2, $3, $4, $5, $6)`,
		event.ID,
		event.Title,
		event.DateTime,
		event.Duration,
		event.Description,
		event.OwnerID,
	)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" {
			switch pgErr.Constraint {
			case "events_pkey":
				return event, storage.NewErrIDNotUnique(event.ID)
			case "events_datetime_owner_id_key":
				return event, storage.NewErrDateBusy(event.OwnerID, event.DateTime)
			}
		}
		return event, err
	}
	if err = tx.Commit(); err != nil {
		return event, err
	}
	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	return storage.Event{}, nil
}
func (s *Storage) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	return nil
}
func (s *Storage) ListEventsDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}
func (s *Storage) ListEventsWeek(ctx context.Context, weekStart time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}
func (s *Storage) ListEventsMonth(ctx context.Context, monthStart time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}
