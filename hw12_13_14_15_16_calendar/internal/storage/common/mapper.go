package storagecommon

import "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"

func ToDomainEvent(e Event) types.Event {
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

func FromDomainEvent(e types.Event) Event {
	return Event{
		ID:           e.ID,
		Title:        e.Title,
		Description:  e.Description,
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		UserID:       e.UserID,
		NotifyBefore: e.NotifyBefore,
	}
}
