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

	storageDB := newSQLStorage()
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

	storageDB := newSQLStorage()
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

func TestStorage_Delete(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC()

	event := storage.Event{
		ID:          "1",
		Title:       "Meeting",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		inputID string
		setup   func(*Storage) storage.EventStorage
		wantErr error
	}{
		{
			name:    "success delete existing event",
			inputID: "1",
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(event)
				return storageDB
			},
			wantErr: nil,
		},
		{
			name:    "fail delete nonexistent event",
			inputID: "2",
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(event)
				return storageDB
			},
			wantErr: storage.ErrEventNotFound,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			s := tt.setup(storageDB)
			defer teardownDB(t, storageDB)

			err := s.Delete(tt.inputID)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				_, err := s.GetByID(tt.inputID)
				require.ErrorIs(t, err, storage.ErrEventNotFound)
			}
		})
	}
}

func TestStorage_GetByID(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC()

	event := storage.Event{
		ID:          "1",
		Title:       "Meeting",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		inputID string
		setup   func(*Storage) storage.EventStorage
		wantErr error
	}{
		{
			name:    "success get existing event",
			inputID: "1",
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(event)
				return storageDB
			},
			wantErr: nil,
		},
		{
			name:    "fail get nonexistent event",
			inputID: "2",
			setup: func(storageDB *Storage) storage.EventStorage {
				_ = storageDB.Create(event)
				return storageDB
			},
			wantErr: storage.ErrEventNotFound,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			s := tt.setup(storageDB)
			defer teardownDB(t, storageDB)

			got, err := s.GetByID(tt.inputID)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, eventToNoTime(event), eventToNoTime(got))
				require.WithinDuration(t, event.StartTime.UTC(), got.StartTime.UTC(), time.Microsecond)
				require.WithinDuration(t, event.EndTime.UTC(), got.EndTime.UTC(), time.Microsecond)
			}
		})
	}
}

func TestStorage_List(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC()

	events := []storage.Event{
		{
			ID:        "1",
			Title:     "Meeting 1",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "2",
			Title:     "Meeting 2",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
			UserID:    "user2",
		},
	}

	tests := []struct {
		name    string
		setup   func(*Storage) storage.EventStorage
		wantLen int
	}{
		{
			name: "list with events",
			setup: func(storageDB *Storage) storage.EventStorage {
				for _, e := range events {
					_ = storageDB.Create(e)
				}
				return storageDB
			},
			wantLen: len(events),
		},
		{
			name: "empty list",
			setup: func(storageDB *Storage) storage.EventStorage {
				return storageDB
			},
			wantLen: 0,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			s := tt.setup(storageDB)
			defer teardownDB(t, storageDB)

			list, err := s.List()
			require.NoError(t, err)
			require.Len(t, list, tt.wantLen)
		})
	}
}

func TestStorage_ListByUser(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC()

	events := []storage.Event{
		{
			ID:        "1",
			Title:     "User1 Event 1",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "2",
			Title:     "User1 Event 2",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "3",
			Title:     "User2 Event",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user2",
		},
	}

	tests := []struct {
		name    string
		userID  string
		setup   func(*Storage) storage.EventStorage
		wantLen int
	}{
		{
			name:   "list user1 events",
			userID: "user1",
			setup: func(storageDB *Storage) storage.EventStorage {
				for _, e := range events {
					_ = storageDB.Create(e)
				}
				return storageDB
			},
			wantLen: 2,
		},
		{
			name:   "list user2 events",
			userID: "user2",
			setup: func(storageDB *Storage) storage.EventStorage {
				for _, e := range events {
					_ = storageDB.Create(e)
				}
				return storageDB
			},
			wantLen: 1,
		},
		{
			name:   "list empty for unknown user",
			userID: "unknown",
			setup: func(storageDB *Storage) storage.EventStorage {
				for _, e := range events {
					_ = storageDB.Create(e)
				}
				return storageDB
			},
			wantLen: 0,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			s := tt.setup(storageDB)
			defer teardownDB(t, storageDB)

			list, err := s.ListByUser(tt.userID)
			require.NoError(t, err)
			require.Len(t, list, tt.wantLen)
		})
	}
}

func TestStorage_ListByUserInRange(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC().Truncate(24 * time.Hour) // нормализуем до начала дня

	events := []storage.Event{
		{
			ID:        "1",
			Title:     "Morning Meeting",
			StartTime: now.Add(9 * time.Hour),
			EndTime:   now.Add(10 * time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "2",
			Title:     "Lunch Break",
			StartTime: now.Add(12 * time.Hour),
			EndTime:   now.Add(13 * time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "3",
			Title:     "Evening Walk",
			StartTime: now.Add(18 * time.Hour),
			EndTime:   now.Add(19 * time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "4",
			Title:     "Another User",
			StartTime: now.Add(10 * time.Hour),
			EndTime:   now.Add(11 * time.Hour),
			UserID:    "user2",
		},
	}

	tests := []struct {
		name    string
		userID  string
		from    time.Time
		to      time.Time
		wantIDs []string
	}{
		{
			name:    "range covers first two events",
			userID:  "user1",
			from:    now.Add(8 * time.Hour),
			to:      now.Add(12*time.Hour + 30*time.Minute),
			wantIDs: []string{"1", "2"},
		},
		{
			name:    "range covers only second event",
			userID:  "user1",
			from:    now.Add(12*time.Hour + 15*time.Minute),
			to:      now.Add(12*time.Hour + 45*time.Minute),
			wantIDs: []string{"2"},
		},
		{
			name:    "range has no events",
			userID:  "user1",
			from:    now.Add(20 * time.Hour),
			to:      now.Add(21 * time.Hour),
			wantIDs: []string{},
		},
		{
			name:    "other user's events not included",
			userID:  "user2",
			from:    now.Add(8 * time.Hour),
			to:      now.Add(12 * time.Hour),
			wantIDs: []string{"4"},
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			s := storageDB
			for _, e := range events {
				_ = s.Create(e)
			}

			list, err := s.ListByUserInRange(tt.userID, tt.from, tt.to)
			require.NoError(t, err)

			var gotIDs []string
			for _, e := range list {
				gotIDs = append(gotIDs, e.ID)
			}

			require.ElementsMatch(t, tt.wantIDs, gotIDs)
			defer teardownDB(t, storageDB)
		})
	}
}

func newSQLStorage() *Storage {
	return New(Config{
		StorageType:    "postgres",
		DSN:            "postgresql://postgres@localhost:5432/calendar_test?sslmode=disable",
		MigrationsPath: filepath.Join(RootDir(), "migrations"),
	})
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
