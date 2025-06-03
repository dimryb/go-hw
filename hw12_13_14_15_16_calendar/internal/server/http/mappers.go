package internalhttp

import (
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

func FromCreateEventRequest(req CreateEventRequest) types.Event {
	return types.Event{
		UserID:       req.UserID,
		Title:        req.Title,
		Description:  req.Description,
		StartTime:    time.Unix(req.StartTime, 0),
		EndTime:      time.Unix(req.EndTime, 0),
		NotifyBefore: int(req.NotifyBefore),
	}
}

func FromUpdateEventRequest(req UpdateEventRequest) types.Event {
	return types.Event{
		ID:           req.ID,
		UserID:       req.UserID,
		Title:        req.Title,
		Description:  req.Description,
		StartTime:    time.Unix(req.StartTime, 0),
		EndTime:      time.Unix(req.EndTime, 0),
		NotifyBefore: int(req.NotifyBefore),
	}
}

func ToEventResponse(event types.Event) EventResponse {
	return EventResponse{
		ID:           event.ID,
		UserID:       event.UserID,
		Title:        event.Title,
		Description:  event.Description,
		StartTime:    event.StartTime.Unix(),
		EndTime:      event.EndTime.Unix(),
		NotifyBefore: int64(event.NotifyBefore),
	}
}

func ToCreateEventRequest(event types.Event) CreateEventRequest {
	return CreateEventRequest{
		UserID:       event.UserID,
		Title:        event.Title,
		Description:  event.Description,
		StartTime:    event.StartTime.Unix(),
		EndTime:      event.EndTime.Unix(),
		NotifyBefore: int64(event.NotifyBefore),
	}
}
