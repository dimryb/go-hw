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
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/mappers"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

//go:generate mockgen -source=app.go -package=mocks -destination=../../mocks/mock_application.go
type Application interface {
	CreateEvent(context.Context, types.Event) (string, error)
	UpdateEvent(context.Context, types.Event) error
	DeleteEvent(context.Context, string) error
	GetEventByID(context.Context, string) (types.Event, error)
	ListEvents(context.Context) ([]types.Event, error)
	ListEventsByUser(context.Context, string) ([]types.Event, error)
	ListEventsByUserInRange(context.Context, string, time.Time, time.Time) ([]types.Event, error)
}

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type Storage interface {
	Create(event storagecommon.Event) (string, error)
	Update(event storagecommon.Event) error
	Delete(id string) error

	GetByID(id string) (storagecommon.Event, error)
	List() ([]storagecommon.Event, error)
	ListByUser(userID string) ([]storagecommon.Event, error)
	ListByUserInRange(userID string, from, to time.Time) ([]storagecommon.Event, error)
}

func Run(configPath string, migrate bool) {
	cfg, err := config.NewCalendarConfig(configPath)
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
		Migration:      migrate,
	})
	if err != nil {
		logg.Fatalf("Failed init storage: %s", err.Error())
	}

	calendar := &App{
		Logger:  logg,
		Storage: storageApp,
	}

	handlers := internalhttp.NewCalendarHandlers(calendar, logg)

	server := internalhttp.NewServer(calendar, logg, internalhttp.ServerConfig{
		Host:              cfg.HTTP.Host,
		Port:              cfg.HTTP.Port,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	},
		handlers,
	)

	if cfg.GRPC.Enable {
		go func() {
			logg.Debugf("gRPC server starting..")
			grpcServer := grpc.NewServer(
				calendar,
				grpc.ServerConfig{
					Port: cfg.GRPC.Port,
				},
				logg,
			)
			if err := grpcServer.Run(); err != nil {
				logg.Fatalf("Failed to start gRPC server: %s", err.Error())
			}
			logg.Debugf("gRPC server started")
		}()
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Errorf("failed to stop http server: %s", err.Error())
		}
	}()

	logg.Infof("calendar is running...")

	if err := server.Start(ctx); err != nil {
		cancel()
		logg.Fatalf("failed to start http server: %s", err.Error())
	}
}

func (a *App) CreateEvent(_ context.Context, event types.Event) (string, error) {
	storEvent := mappers.FromDomainEvent(event)
	id, err := a.Storage.Create(storEvent)
	return id, err
}

func (a *App) UpdateEvent(_ context.Context, event types.Event) error {
	storEvent := mappers.FromDomainEvent(event)
	return a.Storage.Update(storEvent)
}

func (a *App) DeleteEvent(_ context.Context, id string) error {
	return a.Storage.Delete(id)
}

func (a *App) GetEventByID(_ context.Context, id string) (types.Event, error) {
	storEvent, err := a.Storage.GetByID(id)
	if err != nil {
		return types.Event{}, err
	}
	return mappers.ToDomainEvent(storEvent), nil
}

func (a *App) ListEvents(_ context.Context) ([]types.Event, error) {
	storEvents, err := a.Storage.List()
	if err != nil {
		return nil, err
	}

	domainEvents := make([]types.Event, 0, len(storEvents))
	for _, storEvent := range storEvents {
		domainEvents = append(domainEvents, mappers.ToDomainEvent(storEvent))
	}

	return domainEvents, nil
}

func (a *App) ListEventsByUser(_ context.Context, userID string) ([]types.Event, error) {
	storEvents, err := a.Storage.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	domainEvents := make([]types.Event, 0, len(storEvents))
	for _, storEvent := range storEvents {
		domainEvents = append(domainEvents, mappers.ToDomainEvent(storEvent))
	}

	return domainEvents, nil
}

func (a *App) ListEventsByUserInRange(
	_ context.Context,
	userID string,
	from, to time.Time,
) ([]types.Event, error) {
	storEvents, err := a.Storage.ListByUserInRange(userID, from, to)
	if err != nil {
		return nil, err
	}

	domainEvents := make([]types.Event, 0, len(storEvents))
	for _, storEvent := range storEvents {
		domainEvents = append(domainEvents, mappers.ToDomainEvent(storEvent))
	}

	return domainEvents, nil
}
