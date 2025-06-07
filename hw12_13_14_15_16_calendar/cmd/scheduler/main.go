package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/mappers"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

var configFile string

type App struct {
	Storage i.Storage
	Logger  i.Logger
}

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	Run(configFile)
}

func Run(configPath string) {
	cfg, err := config.NewSchedulerConfig(configPath)
	if err != nil {
		log.Fatalf("SchedulerConfig error: %s", err)
	}

	logg := logger.New(cfg.Log.Level)
	logg.Debugf("Scheduler Config: %v", *cfg)

	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User, cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host, cfg.RabbitMQ.Port,
	)
	logg.Debugf("AMQP URL: %s", amqpURL)

	rmqClient, err := rmq.NewClient(amqpURL, cfg.RabbitMQ.Exchange)
	if err != nil {
		logg.Fatalf("Failed to create RMQ client: %v", err)
	}
	defer func() {
		err = rmqClient.Close()
		logg.Fatalf("Failed to close RMQ client: %v", err)
	}()

	var storageApp i.Storage
	storageApp, err = storage.InitStorage(
		storage.Config{
			Type:      cfg.Database.Type,
			DSN:       cfg.Database.DSN,
			Timeout:   cfg.Database.Timeout,
			Migration: false,
		},
	)
	if err != nil {
		logg.Fatalf("Failed to initialize storage: %v", err)
	}

	app := NewApp(storageApp, logg)

	interval := cfg.Scheduler.Interval
	logg.Infof("Scheduler started with interval: %v", interval)

	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			now := time.Now()
			events, err := app.ListEventsDueBefore(ctx, now.Add(interval))
			if err != nil {
				logg.Errorf("Error fetching events: %v", err)
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
				if err := rmqClient.Publish(event.UserID, body); err != nil {
					logg.Errorf("Failed to publish notification for event %s: %v", event.ID, err)
					continue
				}
				logg.Infof("Published notification for event %s", event.ID)
			}
		}
	}
}

func NewApp(storage i.Storage, logger i.Logger) *App {
	return &App{
		Storage: storage,
		Logger:  logger,
	}
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
