package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/proto/calendar"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrEventNotFound   = status.Error(codes.NotFound, "event not found")
	ErrAlreadyExists   = status.Error(codes.AlreadyExists, "event already exists")
	ErrConflictOverlap = status.Error(codes.FailedPrecondition, "event overlaps with existing one")
	ErrInternal        = status.Error(codes.Internal, "internal server error")
)

type Storage interface {
	Create(event storagecommon.Event) error
	GetByID(id string) (storagecommon.Event, error)
	Update(event storagecommon.Event) error
	Delete(id string) error
	List() ([]storagecommon.Event, error)
	ListByUser(userID string) ([]storagecommon.Event, error)
	ListByUserInRange(userID string, from, to time.Time) ([]storagecommon.Event, error)
}

type CalendarService struct {
	calendar.UnimplementedCalendarServiceServer
	storage Storage
}

func NewCalendarService(storage Storage) *CalendarService {
	return &CalendarService{
		storage: storage,
	}
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

func translateError(err error) error {
	switch {
	case errors.Is(err, storagecommon.ErrEventNotFound):
		return ErrEventNotFound
	case errors.Is(err, storagecommon.ErrAlreadyExists):
		return ErrAlreadyExists
	case errors.Is(err, storagecommon.ErrConflictOverlap):
		return ErrConflictOverlap
	default:
		return ErrInternal
	}
}

func (s *CalendarService) CreateEvent(
	_ context.Context,
	event *calendar.Event,
) (*calendar.CreateEventResponse, error) {
	domainEvent := s.protoToDomain(event)
	if err := s.storage.Create(domainEvent); err != nil {
		return nil, translateError(err)
	}
	return &calendar.CreateEventResponse{
		Id:      domainEvent.ID,
		Success: true,
	}, nil
}

func (s *CalendarService) UpdateEvent(
	_ context.Context,
	event *calendar.Event,
) (*calendar.UpdateEventResponse, error) {
	domainEvent := s.protoToDomain(event)
	if err := s.storage.Update(domainEvent); err != nil {
		return nil, translateError(err)
	}
	return &calendar.UpdateEventResponse{Success: true}, nil
}

func (s *CalendarService) DeleteEvent(
	_ context.Context,
	req *calendar.DeleteEventRequest,
) (*calendar.DeleteEventResponse, error) {
	if err := s.storage.Delete(req.Id); err != nil {
		return nil, translateError(err)
	}
	return &calendar.DeleteEventResponse{Success: true}, nil
}

func (s *CalendarService) GetEventByID(
	_ context.Context,
	req *calendar.GetEventByIDRequest,
) (*calendar.GetEventByIDResponse, error) {
	event, err := s.storage.GetByID(req.Id)
	if err != nil {
		return nil, translateError(err)
	}
	return &calendar.GetEventByIDResponse{
		Event: s.domainToProto(event),
	}, nil
}

func (s *CalendarService) ListEvents(
	_ context.Context,
	_ *calendar.ListEventsRequest,
) (*calendar.ListEventsResponse, error) {
	events, err := s.storage.List()
	if err != nil {
		return nil, translateError(err)
	}
	protoEvents := make([]*calendar.Event, 0, len(events))
	for _, e := range events {
		protoEvents = append(protoEvents, s.domainToProto(e))
	}
	return &calendar.ListEventsResponse{Events: protoEvents}, nil
}

func (s *CalendarService) ListEventsByUser(
	_ context.Context,
	req *calendar.ListEventsByUserRequest,
) (*calendar.ListEventsResponse, error) {
	events, err := s.storage.ListByUser(req.UserId)
	if err != nil {
		return nil, translateError(err)
	}
	protoEvents := make([]*calendar.Event, 0, len(events))
	for _, e := range events {
		protoEvents = append(protoEvents, s.domainToProto(e))
	}
	return &calendar.ListEventsResponse{Events: protoEvents}, nil
}

func (s *CalendarService) ListEventsByUserInRange(
	_ context.Context,
	req *calendar.ListEventsByUserInRangeRequest,
) (*calendar.ListEventsResponse, error) {
	from := time.Unix(req.From, 0)
	to := time.Unix(req.To, 0)
	events, err := s.storage.ListByUserInRange(req.UserId, from, to)
	if err != nil {
		return nil, translateError(err)
	}
	protoEvents := make([]*calendar.Event, 0, len(events))
	for _, e := range events {
		protoEvents = append(protoEvents, s.domainToProto(e))
	}
	return &calendar.ListEventsResponse{Events: protoEvents}, nil
}
