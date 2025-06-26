package sender

import (
	"context"
	"encoding/json"
	"time"

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

			if err := s.sendStatus(notif, "delivered"); err != nil {
				s.logger.Errorf("Error sending delivered status: %v", err)
			}
		}
	}
}

func (s *Sender) sendStatus(notification rmq.Notification, status string) error {
	statusMsg := rmq.NotificationStatus{
		NotificationID: notification.ID,
		EventID:        notification.ID, // или передавать отдельно
		UserID:         notification.UserID,
		Status:         status,
		Timestamp:      time.Now(),
	}

	body, err := json.Marshal(statusMsg)
	if err != nil {
		return err
	}

	err = s.rmq.Publish("notification_status", body)
	if err != nil {
		s.logger.Errorf("Failed to publish status for notification %s: %v", notification.ID, err)
		return err
	}

	s.logger.Infof("Published status '%s' for notification %s", status, notification.ID)
	return nil
}
