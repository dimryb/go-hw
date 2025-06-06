package mappers

import (
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

func ToDomainEvent(e storagecommon.Event) types.Event {
	return types.Event{
		ID:           e.ID,
		Title:        e.Title,
		Description:  e.Description,
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		UserID:       e.UserID,
		NotifyBefore: e.NotifyBefore,
	}
}

func ToDomainEvents(events []storagecommon.Event) []types.Event {
	domainEvents := make([]types.Event, 0, len(events))
	for _, e := range events {
		domainEvents = append(domainEvents, ToDomainEvent(e))
	}
	return domainEvents
}

func FromDomainEvent(e types.Event) storagecommon.Event {
	return storagecommon.Event{
		ID:           e.ID,
		Title:        e.Title,
		Description:  e.Description,
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		UserID:       e.UserID,
		NotifyBefore: e.NotifyBefore,
	}
}
