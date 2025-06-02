package sqlstorage

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
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

func eventToNoTime(e storagecommon.Event) EventNoTime {
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
	event := storagecommon.Event{
		Title:       "Meeting",
		StartTime:   now.UTC(),
		EndTime:     now.Add(time.Hour).UTC(),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		input   storagecommon.Event
		setup   func(*Storage) storagecommon.EventStorage
		wantErr error
	}{
		{
			name:  "success create new event",
			input: event,
			setup: func(storageDB *Storage) storagecommon.EventStorage {
				return storageDB
			},
			wantErr: nil,
		},
		{
			name:  "fail event already exists",
			input: event,
			setup: func(storageDB *Storage) storagecommon.EventStorage {
				_, _ = storageDB.Create(event)
				return storageDB
			},
			wantErr: storagecommon.ErrAlreadyExists,
		},
		{
			name: "fail time overlap",
			input: storagecommon.Event{
				Title:     "Another Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user1",
			},
			setup: func(storageDB *Storage) storagecommon.EventStorage {
				_, _ = storageDB.Create(event)
				return storageDB
			},
			wantErr: storagecommon.ErrConflictOverlap,
		},
		{
			name: "success different user same time",
			input: storagecommon.Event{
				Title:     "Another Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user2",
			},
			setup: func(storageDB *Storage) storagecommon.EventStorage {
				_, _ = storageDB.Create(event)
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

			expect := tt.input
			id, err := s.Create(tt.input)
			expect.ID = id

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				got, err := s.GetByID(id)
				require.NoError(t, err)

				require.Equal(t, eventToNoTime(expect), eventToNoTime(got))

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

	baseEvent := storagecommon.Event{
		Title:       "Meeting",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		setup   func(*Storage) (string, error)
		input   func(id string) storagecommon.Event
		wantErr error
	}{
		{
			name: "success update event",
			setup: func(s *Storage) (string, error) {
				return s.Create(baseEvent)
			},
			input: func(id string) storagecommon.Event {
				return baseEvent.WithID(id).With(
					func(e storagecommon.Event) storagecommon.Event {
						e.Title = "Updated Meeting"
						e.StartTime = e.EndTime.Add(-30 * time.Minute)
						e.EndTime = e.StartTime.Add(time.Hour)
						return e
					},
				)
			},
			wantErr: nil,
		},
		{
			name: "fail event not found",
			setup: func(*Storage) (string, error) {
				return "12345678-1234-1234-1234-123456780001", nil
			},
			input:   baseEvent.WithID,
			wantErr: storagecommon.ErrEventNotFound,
		},
		{
			name: "fail time overlap",
			setup: func(s *Storage) (string, error) {
				id, err := s.Create(baseEvent)
				if err != nil {
					return "", err
				}

				_, err = s.Create(storagecommon.Event{
					UserID:      "user1",
					Title:       "Another Event",
					StartTime:   now.Add(time.Hour + 29*time.Minute),
					EndTime:     now.Add(2 * time.Hour),
					Description: "Some other meeting",
				})
				if err != nil {
					return "", err
				}

				return id, nil
			},
			input: func(id string) storagecommon.Event {
				return baseEvent.WithID(id).With(func(e storagecommon.Event) storagecommon.Event {
					e.StartTime = now.Add(30 * time.Minute)
					e.EndTime = now.Add(time.Hour + 30*time.Minute)
					return e
				})
			},
			wantErr: storagecommon.ErrConflictOverlap,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			defer teardownDB(t, storageDB)

			id, err := tt.setup(storageDB)
			require.NoError(t, err)

			input := tt.input(id)

			err = storageDB.Update(input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)

				got, err := storageDB.GetByID(input.ID)
				require.NoError(t, err)

				require.Equal(t, input.UserID, got.UserID)
				require.Equal(t, input.Title, got.Title)
				require.WithinDuration(t, input.StartTime.UTC(), got.StartTime.UTC(), time.Microsecond)
				require.WithinDuration(t, input.EndTime.UTC(), got.EndTime.UTC(), time.Microsecond)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC()

	event := storagecommon.Event{
		Title:       "Meeting",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		inputID *string
		setup   func(*Storage) (storagecommon.EventStorage, string)
		wantErr error
	}{
		{
			name:    "success delete existing event",
			inputID: nil,
			setup: func(storageDB *Storage) (storagecommon.EventStorage, string) {
				id, _ := storageDB.Create(event)
				return storageDB, id
			},
			wantErr: nil,
		},
		{
			name: "fail delete nonexistent event",
			inputID: func() *string {
				id := "12345678-1234-1234-1234-123456780002"
				return &id
			}(),
			setup: func(storageDB *Storage) (storagecommon.EventStorage, string) {
				id, _ := storageDB.Create(event)
				return storageDB, id
			},
			wantErr: storagecommon.ErrEventNotFound,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			s, realID := tt.setup(storageDB)
			defer teardownDB(t, storageDB)

			deleteID := realID
			if tt.inputID != nil {
				deleteID = *tt.inputID
			}
			err := s.Delete(deleteID)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				_, err := s.GetByID(realID)
				require.ErrorIs(t, err, storagecommon.ErrEventNotFound)
			}
		})
	}
}

func TestStorage_GetByID(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC()

	event := storagecommon.Event{
		Title:       "Meeting",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		inputID *string
		setup   func(*Storage) (storagecommon.EventStorage, string)
		wantErr error
	}{
		{
			name:    "success get existing event",
			inputID: nil,
			setup: func(storageDB *Storage) (storagecommon.EventStorage, string) {
				id, _ := storageDB.Create(event)
				return storageDB, id
			},
			wantErr: nil,
		},
		{
			name: "fail get nonexistent event",
			inputID: func() *string {
				id := "12345678-1234-1234-1234-123456780003"
				return &id
			}(),
			setup: func(storageDB *Storage) (storagecommon.EventStorage, string) {
				id, _ := storageDB.Create(event)
				return storageDB, id
			},
			wantErr: storagecommon.ErrEventNotFound,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			s, realID := tt.setup(storageDB)
			defer teardownDB(t, storageDB)

			getID := realID
			if tt.inputID != nil {
				getID = *tt.inputID
			}
			got, err := s.GetByID(getID)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, event.UserID, got.UserID)
				require.Equal(t, event.Title, got.Title)
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

	baseEvents := []storagecommon.Event{
		{
			Title:     "Meeting 1",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user1",
		},
		{
			Title:     "Meeting 2",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
			UserID:    "user2",
		},
	}

	tests := []struct {
		name    string
		setup   func(*Storage) error // добавляем события в БД
		wantLen int
	}{
		{
			name: "list with events",
			setup: func(s *Storage) error {
				for _, e := range baseEvents {
					_, err := s.Create(e)
					if err != nil {
						return err
					}
				}
				return nil
			},
			wantLen: len(baseEvents),
		},
		{
			name: "empty list",
			setup: func(*Storage) error {
				// ничего не создаём
				return nil
			},
			wantLen: 0,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			defer teardownDB(t, storageDB)

			// Шаг 1: подготовка данных
			err := tt.setup(storageDB)
			require.NoError(t, err)

			// Шаг 2: получаем список событий
			list, err := storageDB.List()
			require.NoError(t, err)

			// Шаг 3: проверяем длину
			require.Len(t, list, tt.wantLen)
		})
	}
}

func TestStorage_ListByUser(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC()

	events := []storagecommon.Event{
		{
			Title:     "User1 Event 1",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user1",
		},
		{
			Title:     "User1 Event 2",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
			UserID:    "user1",
		},
		{
			Title:     "User2 Event",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user2",
		},
	}

	tests := []struct {
		name    string
		userID  string
		setup   func(*Storage) ([]string, error)
		wantLen int
	}{
		{
			name:   "list user1 events",
			userID: "user1",
			setup: func(s *Storage) ([]string, error) {
				var ids []string
				for _, e := range events[:2] {
					id, err := s.Create(e)
					if err != nil {
						return nil, err
					}
					ids = append(ids, id)
				}
				return ids, nil
			},
			wantLen: 2,
		},
		{
			name:   "list user2 events",
			userID: "user2",
			setup: func(s *Storage) ([]string, error) {
				e := events[2]
				id, err := s.Create(e)
				if err != nil {
					return nil, err
				}
				return []string{id}, nil
			},
			wantLen: 1,
		},
		{
			name:   "list empty for unknown user",
			userID: "unknown",
			setup: func(s *Storage) ([]string, error) {
				for _, e := range events {
					_, err := s.Create(e)
					if err != nil {
						return nil, err
					}
				}
				return []string{}, nil
			},
			wantLen: 0,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			defer teardownDB(t, storageDB)

			ids, err := tt.setup(storageDB)
			require.NoError(t, err)

			list, err := storageDB.ListByUser(tt.userID)
			require.NoError(t, err)

			require.Len(t, list, tt.wantLen)

			for _, item := range list {
				require.Equal(t, tt.userID, item.UserID)
			}

			if tt.wantLen > 0 {
				require.ElementsMatch(t, ids, extractIDs(list))
			}
		})
	}
}

func TestStorage_ListByUserInRange(t *testing.T) {
	if os.Getenv("TEST_SQL") == "" {
		t.Skip("TEST_SQL not set")
	}

	now := time.Now().UTC().Truncate(24 * time.Hour) // нормализуем до начала дня

	events := []storagecommon.Event{
		{
			Title:     "Morning Meeting",
			StartTime: now.Add(9 * time.Hour),
			EndTime:   now.Add(10 * time.Hour),
			UserID:    "user1",
		},
		{
			Title:     "Lunch Break",
			StartTime: now.Add(12 * time.Hour),
			EndTime:   now.Add(13 * time.Hour),
			UserID:    "user1",
		},
		{
			Title:     "Evening Walk",
			StartTime: now.Add(18 * time.Hour),
			EndTime:   now.Add(19 * time.Hour),
			UserID:    "user1",
		},
		{
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
		wantLen int
	}{
		{
			name:    "range covers first two events",
			userID:  "user1",
			from:    now.Add(8 * time.Hour),
			to:      now.Add(12*time.Hour + 30*time.Minute),
			wantLen: 2,
		},
		{
			name:    "range covers only second event",
			userID:  "user1",
			from:    now.Add(12*time.Hour + 15*time.Minute),
			to:      now.Add(12*time.Hour + 45*time.Minute),
			wantLen: 1,
		},
		{
			name:    "range has no events",
			userID:  "user1",
			from:    now.Add(20 * time.Hour),
			to:      now.Add(21 * time.Hour),
			wantLen: 0,
		},
		{
			name:    "other user's events not included",
			userID:  "user2",
			from:    now.Add(8 * time.Hour),
			to:      now.Add(12 * time.Hour),
			wantLen: 1,
		},
	}

	storageDB := newSQLStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initDB(t, storageDB)
			defer teardownDB(t, storageDB)

			for _, e := range events {
				_, err := storageDB.Create(e)
				require.NoError(t, err)
			}

			list, err := storageDB.ListByUserInRange(tt.userID, tt.from, tt.to)
			require.NoError(t, err)

			require.Len(t, list, tt.wantLen)

			for _, item := range list {
				require.Equal(t, tt.userID, item.UserID)
			}
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

func extractIDs(events []storagecommon.Event) []string {
	ids := make([]string, 0, len(events))
	for _, e := range events {
		ids = append(ids, e.ID)
	}
	return ids
}
