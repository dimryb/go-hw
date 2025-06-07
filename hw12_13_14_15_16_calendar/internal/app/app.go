package app

import (
	"context"
	"time"

	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/mappers"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

type App struct {
	Logger  i.Logger
	Storage i.Storage
}

func NewApp(storage i.Storage) i.Application {
	return &App{Storage: storage}
}

func (a *App) CreateEvent(_ context.Context, event types.Event) (string, error) {
	storEvent := mappers.FromDomainEvent(event)
	id, err := a.Storage.Create(storEvent)
	return id, err
}

func (a *App) UpdateEvent(_ context.Context, event types.Event) error {
	storEvent := mappers.FromDomainEvent(event)
	return a.Storage.Update(storEvent)
}

func (a *App) DeleteEvent(_ context.Context, id string) error {
	return a.Storage.Delete(id)
}

func (a *App) GetEventByID(_ context.Context, id string) (types.Event, error) {
	storEvent, err := a.Storage.GetByID(id)
	if err != nil {
		return types.Event{}, err
	}
	return mappers.ToDomainEvent(storEvent), nil
}

func (a *App) ListEvents(_ context.Context) ([]types.Event, error) {
	storEvents, err := a.Storage.List()
	if err != nil {
		return nil, err
	}

	domainEvents := make([]types.Event, 0, len(storEvents))
	for _, storEvent := range storEvents {
		domainEvents = append(domainEvents, mappers.ToDomainEvent(storEvent))
	}

	return domainEvents, nil
}

func (a *App) ListEventsByUser(_ context.Context, userID string) ([]types.Event, error) {
	storEvents, err := a.Storage.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	domainEvents := make([]types.Event, 0, len(storEvents))
	for _, storEvent := range storEvents {
		domainEvents = append(domainEvents, mappers.ToDomainEvent(storEvent))
	}

	return domainEvents, nil
}

func (a *App) ListEventsByUserInRange(
	_ context.Context,
	userID string,
	from, to time.Time,
) ([]types.Event, error) {
	storEvents, err := a.Storage.ListByUserInRange(userID, from, to)
	if err != nil {
		return nil, err
	}

	domainEvents := make([]types.Event, 0, len(storEvents))
	for _, storEvent := range storEvents {
		domainEvents = append(domainEvents, mappers.ToDomainEvent(storEvent))
	}

	return domainEvents, nil
}

func (a *App) ListEventsDueBefore(_ context.Context, before time.Time) ([]types.Event, error) {
	allEvents, err := a.Storage.List()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dueEvents := make([]types.Event, 0)

	for _, event := range allEvents {
		notifyAt := event.StartTime.Add(-time.Second * time.Duration(event.NotifyBefore))

		if event.StartTime.After(now) && notifyAt.Before(before) && notifyAt.After(now) {
			dueEvents = append(dueEvents, mappers.ToDomainEvent(event))
		}
	}

	return dueEvents, nil
}
