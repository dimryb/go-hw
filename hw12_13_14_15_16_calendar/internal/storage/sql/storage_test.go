package sqlstorage

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
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

	setup := func() storage.EventStorage {
		root := RootDir()
		st := New(Config{
			StorageType:    "postgres",
			DSN:            "postgresql://postgres@localhost:5432/calendar?sslmode=disable",
			MigrationsPath: filepath.Join(root, "migrations"),
		})
		err := st.Connect(context.Background())
		require.NoError(t, err)

		err = st.Migrate()
		require.NoError(t, err)

		return st
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup()
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

func RootDir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(currentFile), "..", "..", "..")
}
