package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/mappers"
	storagecommon "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/mocks"
	pb "github.com/dimryb/go-hw/hw12_13_14_15_calendar/proto/calendar"
	"github.com/golang/mock/gomock" //nolint:depguard
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateEvent(t *testing.T) {
	tests := []struct {
		name        string
		event       *pb.Event
		mockError   error
		expectID    string
		expectError bool
	}{
		{
			name: "Valid Event",
			event: &pb.Event{
				Id:           "event-001",
				UserId:       "user-001",
				Title:        "Birthday",
				Description:  "Party",
				StartTime:    time.Now().Unix() + 3600,
				EndTime:      time.Now().Unix() + 7200,
				NotifyBefore: 3600,
			},
			mockError:   nil,
			expectID:    "event-001",
			expectError: false,
		},
		{
			name: "Already Exists",
			event: &pb.Event{
				Id:           "event-001",
				UserId:       "user-001",
				Title:        "Birthday",
				Description:  "Party",
				StartTime:    time.Now().Unix() + 3600,
				EndTime:      time.Now().Unix() + 7200,
				NotifyBefore: 3600,
			},
			mockError:   status.Error(codes.AlreadyExists, "already exists"),
			expectID:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockApp := mocks.NewMockApplication(ctrl)
			service := &CalendarService{app: mockApp}

			domainEvent := mappers.ProtoToDomain(tt.event)

			mockApp.EXPECT().
				CreateEvent(gomock.Any(), domainEvent).
				Return(tt.expectID, tt.mockError)

			resp, err := service.CreateEvent(context.Background(), tt.event)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectID, resp.Id)
				assert.True(t, resp.Success)
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	tests := []struct {
		name        string
		event       *pb.Event
		mockError   error
		expectError bool
	}{
		{
			name: "Valid Update",
			event: &pb.Event{
				Id:           "event-001",
				UserId:       "user-001",
				Title:        "Updated Title",
				Description:  "Updated Description",
				StartTime:    time.Now().Unix() + 3600,
				EndTime:      time.Now().Unix() + 7200,
				NotifyBefore: 3600,
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "NotFound",
			event: &pb.Event{
				Id:           "event-001",
				UserId:       "user-001",
				Title:        "Updated Title",
				Description:  "Updated Description",
				StartTime:    time.Now().Unix() + 3600,
				EndTime:      time.Now().Unix() + 7200,
				NotifyBefore: 3600,
			},
			mockError:   status.Error(codes.NotFound, "not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockApp := mocks.NewMockApplication(ctrl)
			service := &CalendarService{app: mockApp}

			domainEvent := mappers.ProtoToDomain(tt.event)

			mockApp.EXPECT().
				UpdateEvent(gomock.Any(), domainEvent).
				Return(tt.mockError)

			resp, err := service.UpdateEvent(context.Background(), tt.event)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, resp.Success)
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		mockError   error
		expectError bool
	}{
		{
			name:        "Success",
			id:          "event-001",
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "NotFound",
			id:          "event-002",
			mockError:   status.Error(codes.NotFound, "not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockApp := mocks.NewMockApplication(ctrl)
			service := &CalendarService{app: mockApp}

			mockApp.EXPECT().
				DeleteEvent(gomock.Any(), tt.id).
				Return(tt.mockError)

			req := &pb.DeleteEventRequest{Id: tt.id}
			resp, err := service.DeleteEvent(context.Background(), req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, resp.Success)
			}
		})
	}
}

func TestGetEventByID(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		mockEvent   storagecommon.Event
		mockError   error
		expectError bool
	}{
		{
			name: "Found",
			id:   "event-001",
			mockEvent: storagecommon.Event{
				ID:           "event-001",
				UserID:       "user-001",
				Title:        "Test",
				Description:  "Desc",
				StartTime:    time.Now(),
				EndTime:      time.Now().Add(time.Hour),
				NotifyBefore: 3600,
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "NotFound",
			id:          "event-002",
			mockEvent:   storagecommon.Event{},
			mockError:   status.Error(codes.NotFound, "not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockApp := mocks.NewMockApplication(ctrl)
			service := &CalendarService{app: mockApp}

			domainEvent := mappers.ToDomainEvent(tt.mockEvent)

			mockApp.EXPECT().
				GetEventByID(gomock.Any(), tt.id).
				Return(domainEvent, tt.mockError)

			req := &pb.GetEventByIDRequest{Id: tt.id}
			resp, err := service.GetEventByID(context.Background(), req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp.Event)
				assert.Equal(t, tt.id, resp.Event.Id)
			}
		})
	}
}

func TestListEventsByUser(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		userID      string
		mockEvents  []storagecommon.Event
		mockError   error
		expectCount int
		expectError bool
	}{
		{
			name:   "Has Events",
			userID: "user-001",
			mockEvents: []storagecommon.Event{
				{
					ID:           "event-001",
					UserID:       "user-001",
					Title:        "Meeting",
					StartTime:    now,
					EndTime:      now.Add(time.Hour),
					NotifyBefore: 3600,
				},
			},
			mockError:   nil,
			expectCount: 1,
			expectError: false,
		},
		{
			name:        "No Events",
			userID:      "user-002",
			mockEvents:  []storagecommon.Event{},
			mockError:   nil,
			expectCount: 0,
			expectError: false,
		},
		{
			name:        "Internal Error",
			userID:      "user-003",
			mockEvents:  nil,
			mockError:   status.Error(codes.Internal, "internal error"),
			expectCount: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockApp := mocks.NewMockApplication(ctrl)
			service := &CalendarService{app: mockApp}
			domainEvents := mappers.ToDomainEvents(tt.mockEvents)

			mockApp.EXPECT().
				ListEventsByUser(gomock.Any(), tt.userID).
				Return(domainEvents, tt.mockError)

			req := &pb.ListEventsByUserRequest{UserId: tt.userID}
			resp, err := service.ListEventsByUser(context.Background(), req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, resp.Events, tt.expectCount)
			}
		})
	}
}

func TestListEventsByUserInRange(t *testing.T) {
	now := time.Now()
	from := now
	to := now.Add(24 * time.Hour)

	tests := []struct {
		name        string
		userID      string
		from        int64
		to          int64
		mockEvents  []storagecommon.Event
		mockError   error
		expectCount int
		expectError bool
	}{
		{
			name:   "One Event",
			userID: "user-001",
			from:   from.Unix(),
			to:     to.Unix(),
			mockEvents: []storagecommon.Event{
				{
					ID:           "event-001",
					UserID:       "user-001",
					Title:        "Meeting",
					StartTime:    from.Add(time.Hour),
					EndTime:      from.Add(2 * time.Hour),
					NotifyBefore: 3600,
				},
			},
			mockError:   nil,
			expectCount: 1,
			expectError: false,
		},
		{
			name:        "Empty List",
			userID:      "user-002",
			from:        from.Unix(),
			to:          to.Unix(),
			mockEvents:  []storagecommon.Event{},
			mockError:   nil,
			expectCount: 0,
			expectError: false,
		},
		{
			name:        "ServerError",
			userID:      "user-003",
			from:        from.Unix(),
			to:          to.Unix(),
			mockEvents:  nil,
			mockError:   status.Error(codes.Internal, "server error"),
			expectCount: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockApp := mocks.NewMockApplication(ctrl)
			service := &CalendarService{app: mockApp}
			domainEvents := mappers.ToDomainEvents(tt.mockEvents)

			mockApp.EXPECT().
				ListEventsByUserInRange(gomock.Any(), tt.userID, time.Unix(tt.from, 0), time.Unix(tt.to, 0)).
				Return(domainEvents, tt.mockError)

			req := &pb.ListEventsByUserInRangeRequest{
				UserId: tt.userID,
				From:   tt.from,
				To:     tt.to,
			}

			resp, err := service.ListEventsByUserInRange(context.Background(), req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, resp.Events, tt.expectCount)
			}
		})
	}
}

func TestListEvents(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		mockEvents  []storagecommon.Event
		mockError   error
		expectCount int
		expectError bool
	}{
		{
			name: "Two Events",
			mockEvents: []storagecommon.Event{
				{
					ID:           "event-001",
					UserID:       "user-001",
					Title:        "Event 1",
					StartTime:    now,
					EndTime:      now.Add(time.Hour),
					NotifyBefore: 3600,
				},
				{
					ID:           "event-002",
					UserID:       "user-002",
					Title:        "Event 2",
					StartTime:    now,
					EndTime:      now.Add(time.Hour),
					NotifyBefore: 3600,
				},
			},
			mockError:   nil,
			expectCount: 2,
			expectError: false,
		},
		{
			name:        "Empty List",
			mockEvents:  []storagecommon.Event{},
			mockError:   nil,
			expectCount: 0,
			expectError: false,
		},
		{
			name:        "Internal Error",
			mockEvents:  nil,
			mockError:   status.Error(codes.Internal, "server error"),
			expectCount: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockApp := mocks.NewMockApplication(ctrl)
			service := &CalendarService{app: mockApp}
			domainEvents := mappers.ToDomainEvents(tt.mockEvents)

			mockApp.EXPECT().
				ListEvents(gomock.Any()).
				Return(domainEvents, tt.mockError)

			req := &pb.ListEventsRequest{}
			resp, err := service.ListEvents(context.Background(), req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, resp.Events, tt.expectCount)
			}
		})
	}
}
