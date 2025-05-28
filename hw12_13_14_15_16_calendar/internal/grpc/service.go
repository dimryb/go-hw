package grpc

import (
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/proto/calendar"
)

type CalendarService struct {
	calendar.UnimplementedCalendarServiceServer
}

func NewCalendarService() *CalendarService {
	return &CalendarService{}
}
