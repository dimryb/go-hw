package memorystorage

import (
	"sort"
	"testing"
	"time"

	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_Create(t *testing.T) {
	now := time.Now()

	event := storagecommon.Event{
		ID:          "1",
		Title:       "Meeting",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		input   storagecommon.Event
		setup   func() i.Storage
		wantErr error
	}{
		{
			name:  "success create new event",
			input: event,
			setup: func() i.Storage {
				return New()
			},
			wantErr: nil,
		},
		{
			name:  "fail event already exists",
			input: event,
			setup: func() i.Storage {
				s := New()
				_, _ = s.Create(event)
				return s
			},
			wantErr: storagecommon.ErrAlreadyExists,
		},
		{
			name: "fail time overlap",
			input: storagecommon.Event{
				ID:        "2",
				Title:     "Another Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user1",
			},
			setup: func() i.Storage {
				s := New()
				_, _ = s.Create(event)
				return s
			},
			wantErr: storagecommon.ErrConflictOverlap,
		},
		{
			name: "success different user same time",
			input: storagecommon.Event{
				ID:        "2",
				Title:     "Another Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user2",
			},
			setup: func() i.Storage {
				s := New()
				_, _ = s.Create(event)
				return s
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			_, err := s.Create(tt.input)

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

	baseEvent := storagecommon.Event{
		ID:          "1",
		Title:       "Meeting",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		Description: "Discuss project",
		UserID:      "user1",
	}

	tests := []struct {
		name    string
		input   storagecommon.Event
		setup   func() i.Storage
		wantErr error
	}{
		{
			name:  "success update event",
			input: baseEvent,
			setup: func() i.Storage {
				s := New()
				_, _ = s.Create(baseEvent)
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
			input: storagecommon.Event{
				ID:        "2",
				Title:     "Non-existent",
				StartTime: now,
				EndTime:   now.Add(time.Hour),
				UserID:    "user1",
			},
			setup: func() i.Storage {
				s := New()
				_, _ = s.Create(baseEvent)
				return s
			},
			wantErr: storagecommon.ErrEventNotFound,
		},
		{
			name: "fail time overlap",
			input: storagecommon.Event{
				ID:        "1",
				Title:     "Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user1",
			},
			setup: func() i.Storage {
				s := New()
				_, _ = s.Create(baseEvent)
				_, _ = s.Create(storagecommon.Event{
					ID:        "2",
					Title:     "Another",
					StartTime: now.Add(time.Hour + 29*time.Minute),
					EndTime:   now.Add(2 * time.Hour),
					UserID:    "user1",
				})
				return s
			},
			wantErr: storagecommon.ErrConflictOverlap,
		},
		{
			name: "success update with no overlap",
			input: storagecommon.Event{
				ID:        "1",
				Title:     "Meeting",
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(time.Hour + 30*time.Minute),
				UserID:    "user1",
			},
			setup: func() i.Storage {
				s := New()
				_, _ = s.Create(baseEvent)
				_, _ = s.Create(storagecommon.Event{
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

	event := storagecommon.Event{
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
		setup   func() i.Storage
		wantErr error
	}{
		{
			name:    "success delete existing event",
			inputID: "1",
			setup: func() i.Storage {
				s := New()
				_, _ = s.Create(event)
				return s
			},
			wantErr: nil,
		},
		{
			name:    "fail delete nonexistent event",
			inputID: "2",
			setup: func() i.Storage {
				s := New()
				_, _ = s.Create(event)
				return s
			},
			wantErr: storagecommon.ErrEventNotFound,
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
				require.ErrorIs(t, err, storagecommon.ErrEventNotFound)
			}
		})
	}
}

func TestStorage_DeleteOlder(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		cutoffTime  time.Time
		setupEvents []storagecommon.Event
		expectedIDs []string
	}{
		{
			name:       "delete past events",
			cutoffTime: now,
			setupEvents: []storagecommon.Event{
				{
					ID:          "1",
					Title:       "Past Event",
					StartTime:   now.Add(-2 * time.Hour),
					EndTime:     now.Add(-1 * time.Hour),
					Description: "Should be deleted",
					UserID:      "user1",
				},
				{
					ID:          "2",
					Title:       "Future Event",
					StartTime:   now.Add(1 * time.Hour),
					EndTime:     now.Add(2 * time.Hour),
					Description: "Should stay",
					UserID:      "user1",
				},
			},
			expectedIDs: []string{"2"},
		},
		{
			name:       "no deletion if all events are newer",
			cutoffTime: now.Add(-1 * time.Hour), // на час раньше
			setupEvents: []storagecommon.Event{
				{
					ID:          "3",
					Title:       "Event A",
					StartTime:   now,
					EndTime:     now.Add(1 * time.Hour),
					Description: "Should stay",
					UserID:      "user1",
				},
				{
					ID:          "4",
					Title:       "Event B",
					StartTime:   now.Add(2 * time.Hour),
					EndTime:     now.Add(3 * time.Hour),
					Description: "Should stay",
					UserID:      "user1",
				},
			},
			expectedIDs: []string{"3", "4"},
		},
		{
			name:       "all events should be deleted",
			cutoffTime: now.Add(1 * time.Hour),
			setupEvents: []storagecommon.Event{
				{
					ID:          "5",
					Title:       "Event A",
					StartTime:   now.Add(-3 * time.Hour),
					EndTime:     now.Add(-2 * time.Hour),
					Description: "Deleted",
					UserID:      "user1",
				},
				{
					ID:          "6",
					Title:       "Event B",
					StartTime:   now.Add(-1 * time.Hour),
					EndTime:     now.Add(-30 * time.Minute),
					Description: "Deleted",
					UserID:      "user1",
				},
			},
			expectedIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()

			for _, event := range tt.setupEvents {
				_, _ = s.Create(event)
			}

			err := s.DeleteOlder(tt.cutoffTime)
			require.NoError(t, err)

			actualIDs := make([]string, 0)
			for id := range s.events {
				actualIDs = append(actualIDs, id)
			}

			sort.Strings(actualIDs)
			sort.Strings(tt.expectedIDs)

			assert.Equal(t, tt.expectedIDs, actualIDs)
		})
	}
}

func TestStorage_List(t *testing.T) {
	now := time.Now()

	events := []storagecommon.Event{
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
		setup   func() i.Storage
		wantLen int
	}{
		{
			name: "list with events",
			setup: func() i.Storage {
				s := New()
				for _, e := range events {
					_, _ = s.Create(e)
				}
				return s
			},
			wantLen: len(events),
		},
		{
			name: "empty list",
			setup: func() i.Storage {
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

func TestStorage_ListByUser(t *testing.T) {
	now := time.Now()

	events := []storagecommon.Event{
		{
			ID:        "1",
			Title:     "User1 Event",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "2",
			Title:     "User1 Another",
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
		name      string
		userID    string
		setup     func() i.Storage
		wantCount int
	}{
		{
			name:   "list user1 events",
			userID: "user1",
			setup: func() i.Storage {
				s := New()
				for _, e := range events {
					_, _ = s.Create(e)
				}
				return s
			},
			wantCount: 2,
		},
		{
			name:   "list user2 events",
			userID: "user2",
			setup: func() i.Storage {
				s := New()
				for _, e := range events {
					_, _ = s.Create(e)
				}
				return s
			},
			wantCount: 1,
		},
		{
			name:   "list empty",
			userID: "user3",
			setup: func() i.Storage {
				s := New()
				for _, e := range events {
					_, _ = s.Create(e)
				}
				return s
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			list, err := s.ListByUser(tt.userID)
			require.NoError(t, err)
			require.Len(t, list, tt.wantCount)
		})
	}
}

func TestStorage_ListByUserInRange(t *testing.T) {
	now := time.Now()

	events := []storagecommon.Event{
		{
			ID:        "1",
			Title:     "Morning",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "2",
			Title:     "Noon",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "3",
			Title:     "Evening",
			StartTime: now.Add(6 * time.Hour),
			EndTime:   now.Add(7 * time.Hour),
			UserID:    "user1",
		},
		{
			ID:        "4",
			Title:     "Other User",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user2",
		},
	}

	tests := []struct {
		name      string
		userID    string
		from      time.Time
		to        time.Time
		wantCount int
	}{
		{
			name:      "day range with 2 events",
			userID:    "user1",
			from:      now,
			to:        now.Add(5 * time.Hour),
			wantCount: 2,
		},
		{
			name:      "day range with 1 event",
			userID:    "user1",
			from:      now.Add(5 * time.Hour),
			to:        now.Add(8 * time.Hour),
			wantCount: 1,
		},
		{
			name:      "no events in range",
			userID:    "user1",
			from:      now.Add(10 * time.Hour),
			to:        now.Add(12 * time.Hour),
			wantCount: 0,
		},
		{
			name:      "other user event not included",
			userID:    "user2",
			from:      now,
			to:        now.Add(5 * time.Hour),
			wantCount: 1,
		},
	}

	setup := func() i.Storage {
		s := New()
		for _, e := range events {
			_, _ = s.Create(e)
		}
		return s
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup()
			list, err := s.ListByUserInRange(tt.userID, tt.from, tt.to)
			require.NoError(t, err)
			require.Len(t, list, tt.wantCount)
		})
	}
}
