package app

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
}

type Storage interface {
	Create(event storagecommon.Event) error
	Update(event storagecommon.Event) error
	Delete(id string) error

	GetByID(id string) (storagecommon.Event, error)
	List() ([]storagecommon.Event, error)
	ListByUser(userID string) ([]storagecommon.Event, error)
	ListByUserInRange(userID string, from, to time.Time) ([]storagecommon.Event, error)
}

func Run(configPath string) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logg := logger.New(cfg.Log.Level)

	fmt.Println(cfg.Database)

	var storageApp Storage
	storageApp, err = storage.InitStorage(storage.Config{
		Type:           cfg.Database.Type,
		DSN:            cfg.Database.DSN,
		MigrationsPath: cfg.Database.MigrationsPath,
		Timeout:        cfg.Database.Timeout,
		Migration:      true,
	})
	if err != nil {
		logg.Fatal(err.Error())
	}

	calendar := &App{
		logger:  logg,
		storage: storageApp,
	}

	server := internalhttp.NewServer(logg, calendar, internalhttp.ServerConfig{
		Host:              cfg.Host,
		Port:              cfg.Port,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	})

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
		cancel()
		logg.Fatal("failed to start http server: " + err.Error())
	}
}

func (a *App) CreateEvent(_ context.Context, id, title string) error {
	return a.storage.Create(storagecommon.Event{ID: id, Title: title})
}
