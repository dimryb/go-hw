package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/service"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
)

var (
	configFile string
	migrate    bool
)

func init() {
	flag.StringVar(&configFile, "config", "configs/calendar.yaml", "Path to configuration file")
	flag.BoolVar(&migrate, "migrate", false, "Migrate DB")
}

// @title Receipt Hub API
// @version 1.0
// @description This is a server for Calendar
// @host localhost:8080
// @BasePath /
// .
func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	run(configFile, migrate)
}

func run(configPath string, migrate bool) {
	cfg, err := config.NewCalendarConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logg := logger.New(cfg.Log.Level)

	var storageImpl i.Storage
	storageImpl, err = storage.InitStorage(storage.Config{
		Type:           cfg.Database.Type,
		DSN:            cfg.Database.DSN,
		MigrationsPath: cfg.Database.MigrationsPath,
		Timeout:        cfg.Database.Timeout,
		Migration:      migrate,
	})
	if err != nil {
		logg.Fatalf("Failed to initialize storage: %v", err)
	}

	appl := app.NewApp(storageImpl)
	calendarService := service.NewCalendar(appl, logg, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		logg.Infof("Shutting down...")
		cancel()
	}()

	logg.Infof("Starting calendar service...")
	if err = calendarService.Run(ctx); err != nil {
		logg.Errorf("Calendar service stopped with error: %v", err)
		cancel()
	} else {
		logg.Infof("Calendar service stopped gracefully")
	}
}
