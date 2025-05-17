package sqlstorage

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
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
		wantErr error
	}{
		{
			name:    "success create new event",
			input:   event,
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
			s := setupDB(t, storageDB)
			err := s.Create(tt.input)
			defer teardownDB(t, storageDB)

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

func setupDB(t *testing.T, storageDB *Storage) storage.EventStorage {
	t.Helper()

	err := storageDB.Connect(context.Background())
	require.NoError(t, err)

	err = storageDB.Migrate()
	require.NoError(t, err)

	return storageDB
}

func teardownDB(t *testing.T, storageDB *Storage) {
	t.Helper()

	err := goose.DownTo(storageDB.db.DB, storageDB.migrationsPath, 0)
	require.NoError(t, err)

	err = storageDB.db.Close()
	require.NoError(t, err)
}

func RootDir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(currentFile), "..", "..", "..")
}
