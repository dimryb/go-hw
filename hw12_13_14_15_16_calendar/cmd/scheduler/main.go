package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/service/scheduler"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	run(configFile)
}

func run(configPath string) {
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

	application := app.NewApp(storageApp, logg)
	schedulerService := scheduler.NewScheduler(application, rmqClient, logg, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		logg.Infof("Shutting down scheduler...")
		if err = rmqClient.Close(); err != nil {
			logg.Errorf("Failed to close RMQ client: %v", err)
		}
		cancel()
	}()

	logg.Infof("Starting scheduler service...")
	if err = schedulerService.Run(ctx); err != nil {
		logg.Errorf("Scheduler service stopped with error: %v", err)
		cancel()
	} else {
		logg.Infof("Scheduler service stopped gracefully")
	}
}
