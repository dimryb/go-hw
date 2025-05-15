package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

type Config struct {
	StorageType    string
	DSN            string
	MigrationsPath string
}

type Storage struct {
	storageType    string
	dsn            string
	migrationsPath string
	db             *sqlx.DB
}

func New(cfg Config) *Storage {
	return &Storage{
		storageType:    cfg.StorageType,
		dsn:            cfg.DSN,
		migrationsPath: cfg.MigrationsPath,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, s.storageType, s.dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", s.storageType, err)
	}

	s.db = db

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping %s: %w", s.storageType, err)
	}

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("failed to close DB connection: %w", err)
		}
		s.db = nil
	}
	return nil
}

func (s *Storage) Migrate() error {
	if err := goose.SetDialect(s.storageType); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(nil, s.migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (s *Storage) Create(event storage.Event) error {
	_, err := s.db.NamedExec(`
        INSERT INTO events (
            id, title, start_time, end_time, description, user_id, notify_before
        ) VALUES (
            :id, :title, :start_time, :end_time, :description, :user_id, :notify_before
        )
    `, event)
	return err
}

func (s *Storage) Update(event storage.Event) error {
	res, err := s.db.NamedExec(`
        UPDATE events SET
            title = :title,
            start_time = :start_time,
            end_time = :end_time,
            description = :description,
            user_id = :user_id,
            notify_before = :notify_before
        WHERE id = :id
    `, event)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrEventNotFound
	}
	return nil
}

func (s *Storage) Delete(id string) error {
	res, err := s.db.Exec("DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrEventNotFound
	}
	return nil
}

func (s *Storage) GetByID(id string) (storage.Event, error) {
	var event storage.Event
	err := s.db.Get(&event, "SELECT * FROM events WHERE id = $1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.Event{}, storage.ErrEventNotFound
	}
	return event, err
}

func (s *Storage) List() ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.Select(&events, "SELECT * FROM events")
	return events, err
}

func (s *Storage) ListByUser(userID string) ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.Select(&events, "SELECT * FROM events WHERE user_id = $1", userID)
	return events, err
}

func (s *Storage) ListByUserInRange(userID string, from, to time.Time) ([]storage.Event, error) {
	var events []storage.Event
	query := `
        SELECT * FROM events 
        WHERE user_id = $1
        AND NOT (end_time <= $2 OR start_time >= $3)
    `
	err := s.db.Select(&events, query, userID, from, to)
	return events, err
}
