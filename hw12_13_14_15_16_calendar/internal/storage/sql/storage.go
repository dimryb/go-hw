package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/jmoiron/sqlx"     //nolint:depguard
	_ "github.com/lib/pq"         //nolint:depguard
	"github.com/pressly/goose/v3" //nolint:depguard
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
	if s.db == nil {
		return fmt.Errorf("database connection is not established")
	}

	if err := goose.SetDialect(s.storageType); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(s.db.DB, s.migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (s *Storage) Create(event storagecommon.Event) (string, error) {
	duplicate, err := s.isDuplicate(event)
	if err != nil {
		return "", err
	}
	if duplicate {
		return "", storagecommon.ErrAlreadyExists
	}

	overlap, err := s.isOverlapping(event)
	if err != nil {
		return "", fmt.Errorf("checking overlapping events: %w", err)
	}
	if overlap {
		return "", storagecommon.ErrConflictOverlap
	}

	const query = `
	   INSERT INTO events (
	       user_id, title, start_time, end_time, description, notify_before
	   ) VALUES (
	       :user_id, :title, :start_time, :end_time, :description, :notify_before
	   )
	   RETURNING id`

	var newID string
	namedQuery, err := s.db.PrepareNamed(query)
	if err != nil {
		return "", fmt.Errorf("failed to prepare named query: %w", err)
	}

	err = namedQuery.Get(&newID, event)
	if err != nil {
		return "", fmt.Errorf("failed to create event: %w", err)
	}

	return newID, nil
}

func (s *Storage) Update(event storagecommon.Event) error {
	existing, err := s.GetByID(event.ID)
	if err != nil {
		return err
	}

	if existing.UserID == event.UserID {
		overlap, err := s.isOverlapping(event)
		if err != nil {
			return fmt.Errorf("checking overlapping events: %w", err)
		}
		if overlap {
			return storagecommon.ErrConflictOverlap
		}
	}

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
		return fmt.Errorf("failed to update event: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storagecommon.ErrEventNotFound
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
		return storagecommon.ErrEventNotFound
	}
	return nil
}

func (s *Storage) DeleteOlder(t time.Time) error {
	_, err := s.db.Exec("DELETE FROM events WHERE end_time < $1", t)
	return err
}

func (s *Storage) GetByID(id string) (storagecommon.Event, error) {
	if id == "" {
		return storagecommon.Event{}, storagecommon.ErrEventNotFound
	}

	var event storagecommon.Event
	err := s.db.Get(&event, "SELECT * FROM events WHERE id = $1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return storagecommon.Event{}, storagecommon.ErrEventNotFound
	}
	return event, err
}

func (s *Storage) List() ([]storagecommon.Event, error) {
	var events []storagecommon.Event
	err := s.db.Select(&events, "SELECT * FROM events")
	return events, err
}

func (s *Storage) ListByUser(userID string) ([]storagecommon.Event, error) {
	var events []storagecommon.Event
	err := s.db.Select(&events, "SELECT * FROM events WHERE user_id = $1", userID)
	return events, err
}

func (s *Storage) ListByUserInRange(userID string, from, to time.Time) ([]storagecommon.Event, error) {
	var events []storagecommon.Event
	query := `
        SELECT * FROM events 
        WHERE user_id = $1
        AND NOT (end_time <= $2 OR start_time >= $3)
    `
	err := s.db.Select(&events, query, userID, from, to)
	return events, err
}

func (s *Storage) isOverlapping(event storagecommon.Event) (bool, error) {
	var err error
	var exists bool
	if event.ID == "" {
		query := `
            SELECT EXISTS (
                SELECT 1 FROM events 
                WHERE user_id = $1
                  AND end_time > $2
                  AND start_time < $3
            )`
		err = s.db.Get(&exists, query,
			event.UserID,
			event.StartTime,
			event.EndTime,
		)
	} else {
		query := `
            SELECT EXISTS (
                SELECT 1 FROM events 
                WHERE user_id = $1
                  AND end_time > $2
                  AND start_time < $3
                  AND id != $4
            )`
		err = s.db.Get(&exists, query,
			event.UserID,
			event.StartTime,
			event.EndTime,
			event.ID,
		)
	}
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *Storage) isDuplicate(event storagecommon.Event) (bool, error) {
	const query = `
        SELECT EXISTS (
            SELECT 1 FROM events
            WHERE user_id = :user_id
              AND title = :title
              AND start_time = :start_time
              AND end_time = :end_time
              AND description = :description
              AND notify_before = :notify_before
        )`

	var exists bool
	namedQuery, err := s.db.PrepareNamed(query)
	if err != nil {
		return false, err
	}

	err = namedQuery.Get(&exists, event)
	if err != nil {
		return false, fmt.Errorf("failed to check duplicate: %w", err)
	}

	return exists, nil
}
