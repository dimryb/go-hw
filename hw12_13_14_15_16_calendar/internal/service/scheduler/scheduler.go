package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/rmq"
)

type Scheduler struct {
	app    i.Application
	rmq    i.RmqClient
	logger i.Logger
	cfg    *config.SchedulerConfig
}

func NewScheduler(app i.Application, rmq i.RmqClient, logger i.Logger, cfg *config.SchedulerConfig) *Scheduler {
	return &Scheduler{
		app:    app,
		rmq:    rmq,
		logger: logger,
		cfg:    cfg,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	s.logger.Infof("Scheduler started with interval: %v", s.cfg.Interval)

	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			now := time.Now()
			events, err := s.app.ListEventsDueBefore(ctx, now.Add(s.cfg.Interval))
			if err != nil {
				s.logger.Errorf("Error fetching events: %v", err)
				continue
			}

			for _, event := range events {
				dto := rmq.Notification{
					ID:          event.ID,
					Title:       event.Title,
					Description: event.Description,
					UserID:      event.UserID,
					Time:        event.StartTime.Format(time.RFC3339),
					NotifyAt:    now.Format(time.RFC3339),
				}

				body, _ := json.Marshal(dto)
				if err := s.rmq.Publish(event.UserID, body); err != nil {
					s.logger.Errorf("Failed to publish notification for event %s: %v", event.ID, err)
					continue
				}
				s.logger.Infof("Published notification for event %s", event.ID)
			}

			if err := s.app.DeleteOlderThan(ctx, time.Now().Add(-s.cfg.Scheduler.RetentionPeriod)); err != nil {
				s.logger.Warnf("Failed to delete old events: %v", err)
			}
		}
	}
}
