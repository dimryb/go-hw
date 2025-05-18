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
	memorystorage "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/sql"
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
	Create(event storage.Event) error
	Update(event storage.Event) error
	Delete(id string) error

	GetByID(id string) (storage.Event, error)
	List() ([]storage.Event, error)
	ListByUser(userID string) ([]storage.Event, error)
	ListByUserInRange(userID string, from, to time.Time) ([]storage.Event, error)
}

func Run(configPath string) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logg := logger.New(cfg.Log.Level)

	fmt.Println(cfg.Database)

	var storageApp Storage

	switch cfg.Database.Type {
	case "memory":
		storageApp = memorystorage.New()
	case "postgres":
		sqlStorage := sqlstorage.New(sqlstorage.Config{
			StorageType:    cfg.Database.Type,
			DSN:            cfg.Database.DSN,
			MigrationsPath: cfg.Database.MigrationsPath,
		})

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.Timeout)
		defer cancel()

		if err := sqlStorage.Connect(ctx); err != nil {
			logg.Fatal("failed to connect to database: " + err.Error())
		}

		if err := sqlStorage.Migrate(); err != nil {
			logg.Fatal("failed to apply migrations: " + err.Error())
		}

		storageApp = sqlStorage
	default:
		logg.Fatal("unknown storage type: " + cfg.Database.Type)
	}

	calendar := &App{
		logger:  logg,
		storage: storageApp,
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
		cancel()
		logg.Fatal("failed to start http server: " + err.Error())
	}
}

func (a *App) CreateEvent(_ context.Context, id, title string) error {
	return a.storage.Create(storage.Event{ID: id, Title: title})
}
