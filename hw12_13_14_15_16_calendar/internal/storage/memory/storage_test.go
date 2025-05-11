package memorystorage

import (
	"testing"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage_Create(t *testing.T) {
	now := time.Now()

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
		input   storage.Event
		setup   func() storage.EventStorage
		wantErr error
	}{
		{
			name:  "success create new event",
			input: event,
			setup: func() storage.EventStorage {
				return New()
			},
			wantErr: nil,
		},
		{
			name:  "fail event already exists",
			input: event,
			setup: func() storage.EventStorage {
				s := New()
				_ = s.Create(event)
				return s
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
			setup: func() storage.EventStorage {
				s := New()
				_ = s.Create(event)
				return s
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
			setup: func() storage.EventStorage {
				s := New()
				_ = s.Create(event)
				return s
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			err := s.Create(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				got, err := s.GetByID(tt.input.ID)
				require.NoError(t, err)
				require.Equal(t, tt.input, got)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	now := time.Now()

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
		setup   func() storage.EventStorage
		wantErr error
	}{
		{
			name:  "success update event",
			input: baseEvent,
			setup: func() storage.EventStorage {
				s := New()
				_ = s.Create(baseEvent)
				updated := baseEvent
				updated.Title = "Updated Meeting"
				updated.StartTime = now.Add(2 * time.Hour)
				updated.EndTime = now.Add(3 * time.Hour)
				s.events[baseEvent.ID] = updated
				return s
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
			setup: func() storage.EventStorage {
				s := New()
				_ = s.Create(baseEvent)
				return s
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
			setup: func() storage.EventStorage {
				s := New()
				_ = s.Create(baseEvent)
				_ = s.Create(storage.Event{
					ID:        "2",
					Title:     "Another",
					StartTime: now.Add(90 * time.Minute),
					EndTime:   now.Add(2 * time.Hour),
					UserID:    "user1",
				})
				return s
			},
			wantErr: storage.ErrConflictOverlap,
		},
		{
			name: "no time overlap",
			input: storage.Event{
				ID:        "1",
				Title:     "Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 29*time.Minute),
				UserID:    "user1",
			},
			setup: func() storage.EventStorage {
				s := New()
				_ = s.Create(baseEvent)
				_ = s.Create(storage.Event{
					ID:        "2",
					Title:     "Another",
					StartTime: now.Add(90 * time.Minute),
					EndTime:   now.Add(2 * time.Hour),
					UserID:    "user1",
				})
				return s
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			err := s.Update(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				got, err := s.GetByID(tt.input.ID)
				require.NoError(t, err)
				require.Equal(t, tt.input, got)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	now := time.Now()

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
		setup   func() storage.EventStorage
		wantErr error
	}{
		{
			name:    "success delete existing event",
			inputID: "1",
			setup: func() storage.EventStorage {
				s := New()
				_ = s.Create(event)
				return s
			},
			wantErr: nil,
		},
		{
			name:    "fail delete nonexistent event",
			inputID: "2",
			setup: func() storage.EventStorage {
				s := New()
				_ = s.Create(event)
				return s
			},
			wantErr: storage.ErrEventNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
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

func TestStorage_List(t *testing.T) {
	now := time.Now()

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
		setup   func() storage.EventStorage
		wantLen int
	}{
		{
			name: "list with events",
			setup: func() storage.EventStorage {
				s := New()
				for _, e := range events {
					_ = s.Create(e)
				}
				return s
			},
			wantLen: len(events),
		},
		{
			name: "empty list",
			setup: func() storage.EventStorage {
				return New()
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			list, err := s.List()
			require.NoError(t, err)
			require.Len(t, list, tt.wantLen)
		})
	}
}
