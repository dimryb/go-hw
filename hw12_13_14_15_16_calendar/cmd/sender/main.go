package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/service/sender"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/sender.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	run(configFile)
}

func run(configPath string) {
	cfg, err := config.NewSenderConfig(configPath)
	if err != nil {
		log.Fatalf("Sender config error: %s", err)
	}

	logg := logger.New(cfg.Log.Level)
	logg.Debugf("Sender Config: %v", *cfg)

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

	senderService := sender.NewSender(rmqClient, logg, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		logg.Infof("Shutting down sender...")
		cancel()
	}()

	logg.Infof("Starting sender service...")
	if err = senderService.Run(ctx); err != nil {
		logg.Errorf("Sender service stopped with error: %v", err)
		cancel()
	} else {
		logg.Infof("Sender service stopped gracefully")
	}
}
