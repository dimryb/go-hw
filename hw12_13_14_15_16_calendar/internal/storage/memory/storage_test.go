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
