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
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
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

func Run(configPath string, migrate bool) {
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
		Migration:      migrate,
	})
	if err != nil {
		logg.Fatalf("Failed init storage: %s", err.Error())
	}

	calendar := &App{
		logger:  logg,
		storage: storageApp,
	}

	server := internalhttp.NewServer(logg, calendar, internalhttp.ServerConfig{
		Host:              cfg.HTTP.Host,
		Port:              cfg.HTTP.Port,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	})

	if cfg.GRPC.Enable {
		go func() {
			logg.Debugf("gRPC server starting..")
			grpcServer := grpc.NewServer(
				grpc.ServerConfig{
					Port: cfg.GRPC.Port,
				},
				logg,
				storageApp,
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

func (a *App) CreateEvent(_ context.Context, event types.Event) error {
	storEvent := storagecommon.FromDomainEvent(event)
	return a.storage.Create(storEvent)
}
