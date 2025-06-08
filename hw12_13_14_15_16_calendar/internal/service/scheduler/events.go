package scheduler

import (
	"context"
	"time"

	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

func ListEventsDueBefore(ctx context.Context, app i.Application, before time.Time) ([]types.Event, error) {
	allEvents, err := app.ListEvents(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dueEvents := make([]types.Event, 0)

	for _, event := range allEvents {
		notifyAt := event.StartTime.Add(-time.Second * time.Duration(event.NotifyBefore))

		if event.StartTime.After(now) && notifyAt.Before(before) && notifyAt.After(now) {
			dueEvents = append(dueEvents, event)
		}
	}

	return dueEvents, nil
}
