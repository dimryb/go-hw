package sender

import (
	"context"
	"encoding/json"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/rmq"
)

type Sender struct {
	rmq    i.RmqClient
	logger i.Logger
	cfg    *config.SenderConfig
}

func NewSender(rmq i.RmqClient, logger i.Logger, cfg *config.SenderConfig) *Sender {
	return &Sender{
		rmq:    rmq,
		logger: logger,
		cfg:    cfg,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	s.logger.Infof("Sender started")

	msgChan, err := s.rmq.Consume(s.cfg.QueueName)
	if err != nil {
		s.logger.Errorf("Failed to consume from queue: %v", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("Stopping sender...")
			return ctx.Err()
		case body := <-msgChan:
			var notif rmq.Notification
			if err := json.Unmarshal(body, &notif); err != nil {
				s.logger.Errorf("Failed to unmarshal notification: %v", err)
				continue
			}
			s.logger.Infof("Received notification: %+v", notif)
		}
	}
}
