package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/mappers"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
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
	Create(event types.Event) error
	GetByID(id string) (types.Event, error)
	Update(event types.Event) error
	Delete(id string) error
	List() ([]types.Event, error)
	ListByUser(userID string) ([]types.Event, error)
	ListByUserInRange(userID string, from, to time.Time) ([]types.Event, error)
}

type CalendarService struct {
	calendar.UnimplementedCalendarServiceServer
	app Application
}

func NewCalendarService(app Application) *CalendarService {
	return &CalendarService{
		app: app,
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
	ctx context.Context,
	event *calendar.Event,
) (*calendar.CreateEventResponse, error) {
	domainEvent := mappers.ProtoToDomain(event)
	id, err := s.app.CreateEvent(ctx, domainEvent)
	if err != nil {
		return nil, translateError(err)
	}
	return &calendar.CreateEventResponse{
		Id:      id,
		Success: true,
	}, nil
}

func (s *CalendarService) UpdateEvent(
	ctx context.Context,
	event *calendar.Event,
) (*calendar.UpdateEventResponse, error) {
	domainEvent := mappers.ProtoToDomain(event)
	if err := s.app.UpdateEvent(ctx, domainEvent); err != nil {
		return nil, translateError(err)
	}
	return &calendar.UpdateEventResponse{Success: true}, nil
}

func (s *CalendarService) DeleteEvent(
	ctx context.Context,
	req *calendar.DeleteEventRequest,
) (*calendar.DeleteEventResponse, error) {
	if err := s.app.DeleteEvent(ctx, req.Id); err != nil {
		return nil, translateError(err)
	}
	return &calendar.DeleteEventResponse{Success: true}, nil
}

func (s *CalendarService) GetEventByID(
	ctx context.Context,
	req *calendar.GetEventByIDRequest,
) (*calendar.GetEventByIDResponse, error) {
	event, err := s.app.GetEventByID(ctx, req.Id)
	if err != nil {
		return nil, translateError(err)
	}
	return &calendar.GetEventByIDResponse{
		Event: mappers.DomainToProto(event),
	}, nil
}

func (s *CalendarService) ListEvents(
	ctx context.Context,
	_ *calendar.ListEventsRequest,
) (*calendar.ListEventsResponse, error) {
	events, err := s.app.ListEvents(ctx)
	if err != nil {
		return nil, translateError(err)
	}
	protoEvents := make([]*calendar.Event, 0, len(events))
	for _, e := range events {
		protoEvents = append(protoEvents, mappers.DomainToProto(e))
	}
	return &calendar.ListEventsResponse{Events: protoEvents}, nil
}

func (s *CalendarService) ListEventsByUser(
	ctx context.Context,
	req *calendar.ListEventsByUserRequest,
) (*calendar.ListEventsResponse, error) {
	events, err := s.app.ListEventsByUser(ctx, req.UserId)
	if err != nil {
		return nil, translateError(err)
	}
	protoEvents := make([]*calendar.Event, 0, len(events))
	for _, e := range events {
		protoEvents = append(protoEvents, mappers.DomainToProto(e))
	}
	return &calendar.ListEventsResponse{Events: protoEvents}, nil
}

func (s *CalendarService) ListEventsByUserInRange(
	ctx context.Context,
	req *calendar.ListEventsByUserInRangeRequest,
) (*calendar.ListEventsResponse, error) {
	from := time.Unix(req.From, 0)
	to := time.Unix(req.To, 0)
	events, err := s.app.ListEventsByUserInRange(ctx, req.UserId, from, to)
	if err != nil {
		return nil, translateError(err)
	}
	protoEvents := make([]*calendar.Event, 0, len(events))
	for _, e := range events {
		protoEvents = append(protoEvents, mappers.DomainToProto(e))
	}
	return &calendar.ListEventsResponse{Events: protoEvents}, nil
}
