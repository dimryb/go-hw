package sqlstorage

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/pressly/goose/v3"         //nolint:depguard
	"github.com/stretchr/testify/require" //nolint:depguard
)

type EventNoTime struct {
	ID           string
	Title        string
	Description  string
	UserID       string
	NotifyBefore int
}

func eventToNoTime(e storage.Event) EventNoTime {
	return EventNoTime{
		ID:           e.ID,
		Title:        e.Title,
		Description:  e.Description,
		UserID:       e.UserID,
		NotifyBefore: e.NotifyBefore,
	}
}

func TestStorage_Create(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now()
	event := storage.Event{
		ID:          "1",
		Title:       "Meeting",
		StartTime:   now.UTC(),
		EndTime:     now.Add(time.Hour).UTC(),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		input   storage.Event
		setup   func(*Storage) storage.EventStorage
		wantErr error
	}{
		{
			name:  "success create new event",
			input: event,
			setup: func(storageDB *Storage) storage.EventStorage {
				return storageDB
			},
			wantErr: nil,
		},
		{
			name:  "fail event already exists",
			input: event,
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(event)
				return storageDB
			},
			wantErr: storage.ErrAlreadyExists,
		},
		{
			name: "fail time overlap",
			input: storage.Event{
				ID:        "2",
				Title:     "Another Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user1",
			},
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(event)
				return storageDB
			},
			wantErr: storage.ErrConflictOverlap,
		},
		{
			name: "success different user same time",
			input: storage.Event{
				ID:        "2",
				Title:     "Another Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user2",
			},
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(event)
				return storageDB
			},
			wantErr: nil,
		},
	}

	storageDB := New(Config{
		StorageType:    "postgres",
		DSN:            "postgresql://postgres@localhost:5432/calendar_test?sslmode=disable",
		MigrationsPath: filepath.Join(RootDir(), "migrations"),
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			s := tt.setup(storageDB)
			defer teardownDB(t, storageDB)
			err := s.Create(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				got, err := s.GetByID(tt.input.ID)
				require.NoError(t, err)

				require.Equal(t, eventToNoTime(tt.input), eventToNoTime(got))

				require.WithinDuration(t, tt.input.StartTime.UTC(), got.StartTime.UTC(), time.Microsecond)
				require.WithinDuration(t, tt.input.EndTime.UTC(), got.EndTime.UTC(), time.Microsecond)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC()

	baseEvent := storage.Event{
		ID:          "1",
		Title:       "Meeting",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		input   storage.Event
		setup   func(*Storage) storage.EventStorage
		wantErr error
	}{
		{
			name:  "success update event",
			input: baseEvent,
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(baseEvent)

				updated := baseEvent
				updated.Title = "Updated Meeting"
				updated.StartTime = now.Add(2 * time.Hour)
				updated.EndTime = now.Add(3 * time.Hour)
				storageDB.Update(updated)

				return storageDB
			},
			wantErr: nil,
		},
		{
			name: "fail event not found",
			input: storage.Event{
				ID:        "2",
				Title:     "Non-existent",
				StartTime: now,
				EndTime:   now.Add(time.Hour),
				UserID:    "user1",
			},
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(baseEvent)
				return storageDB
			},
			wantErr: storage.ErrEventNotFound,
		},
		{
			name: "fail time overlap",
			input: storage.Event{
				ID:        "1",
				Title:     "Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user1",
			},
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(baseEvent)
				_ = storageDB.Create(storage.Event{
					ID:        "2",
					Title:     "Another",
					StartTime: now.Add(time.Hour + 29*time.Minute),
					EndTime:   now.Add(2 * time.Hour),
					UserID:    "user1",
				})
				return storageDB
			},
			wantErr: storage.ErrConflictOverlap,
		},
		{
			name: "success update with no overlap",
			input: storage.Event{
				ID:        "1",
				Title:     "Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user1",
			},
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(baseEvent)
				_ = storageDB.Create(storage.Event{
					ID:        "2",
					Title:     "Another",
					StartTime: now.Add(90 * time.Minute),
					EndTime:   now.Add(2 * time.Hour),
					UserID:    "user1",
				})
				return storageDB
			},
			wantErr: nil,
		},
	}

	// Инициализируем хранилище
	storageDB := New(Config{
		StorageType:    "postgres",
		DSN:            "postgresql://postgres@localhost:5432/calendar_test?sslmode=disable",
		MigrationsPath: filepath.Join(RootDir(), "migrations"),
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			s := tt.setup(storageDB)
			defer teardownDB(t, storageDB)

			err := s.Update(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)

				got, err := s.GetByID(tt.input.ID)
				require.NoError(t, err)

				require.Equal(t, eventToNoTime(tt.input), eventToNoTime(got))
				require.WithinDuration(t, tt.input.StartTime.UTC(), got.StartTime.UTC(), time.Microsecond)
				require.WithinDuration(t, tt.input.EndTime.UTC(), got.EndTime.UTC(), time.Microsecond)
			}
		})
	}
}

func initDB(t *testing.T, storageDB *Storage) {
	t.Helper()

	err := storageDB.Connect(context.Background())
	require.NoError(t, err)

	err = storageDB.Migrate()
	require.NoError(t, err)
}

func teardownDB(t *testing.T, storageDB *Storage) {
	t.Helper()

	err := goose.DownTo(storageDB.db.DB, storageDB.migrationsPath, 0)
	require.NoError(t, err)

	err = storageDB.db.Close()
	require.NoError(t, err)
}

func RootDir() string {
	_, currentFile, _, _ := runtime.Caller(0) //nolint:dogsled
	return filepath.Join(filepath.Dir(currentFile), "..", "..", "..")
}
