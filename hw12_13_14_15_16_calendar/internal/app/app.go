package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func Run(configPath string) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logg := logger.New(cfg.Log.Level)

	storage := memorystorage.New()
	calendar := &App{
		logger:  logg,
		storage: storage,
	}

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
