package grpc

import (
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/proto/calendar"
)

type CalendarService struct {
	calendar.UnimplementedCalendarServiceServer
}

func NewCalendarService() *CalendarService {
	return &CalendarService{}
}

func (s *CalendarService) protoToDomain(event *calendar.Event) storagecommon.Event {
	return storagecommon.Event{
		ID:           event.Id,
		UserID:       event.UserId,
		Title:        event.Title,
		Description:  event.Description,
		StartTime:    time.Unix(event.StartTime, 0),
		EndTime:      time.Unix(event.EndTime, 0),
		NotifyBefore: int(event.NotifyBefore),
	}
}

func (s *CalendarService) domainToProto(event storagecommon.Event) *calendar.Event {
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
