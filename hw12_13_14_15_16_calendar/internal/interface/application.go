package interfaces

import (
	"context"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

//go:generate mockgen -source=app.go -package=mocks -destination=../../mocks/mock_application.go
type Application interface {
	CreateEvent(context.Context, types.Event) (string, error)
	UpdateEvent(context.Context, types.Event) error
	DeleteEvent(context.Context, string) error
	GetEventByID(context.Context, string) (types.Event, error)
	ListEvents(context.Context) ([]types.Event, error)
	ListEventsByUser(context.Context, string) ([]types.Event, error)
	ListEventsByUserInRange(context.Context, string, time.Time, time.Time) ([]types.Event, error)

	ListEventsDueBefore(_ context.Context, before time.Time) ([]types.Event, error)
}
