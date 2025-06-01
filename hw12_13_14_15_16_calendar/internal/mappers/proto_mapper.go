package mappers

import (
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/proto/calendar"
)

func ProtoToDomain(event *calendar.Event) types.Event {
	return types.Event{
		ID:           event.Id,
		UserID:       event.UserId,
		Title:        event.Title,
		Description:  event.Description,
		StartTime:    time.Unix(event.StartTime, 0),
		EndTime:      time.Unix(event.EndTime, 0),
		NotifyBefore: int(event.NotifyBefore),
	}
}

func DomainToProto(event types.Event) *calendar.Event {
	return &calendar.Event{
		Id:           event.ID,
		UserId:       event.UserID,
		Title:        event.Title,
		Description:  event.Description,
		StartTime:    event.StartTime.Unix(),
		EndTime:      event.EndTime.Unix(),
		NotifyBefore: int64(event.NotifyBefore),
	}
}
