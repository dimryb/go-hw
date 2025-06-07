package tests

import (
	"context"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/app"
	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage"
	storagecommon "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
)

type TestAppForCalendar struct {
	App     *app.App
	Server  *internalhttp.Server
	Storage storagecommon.EventStorage
	Logger  i.Logger
}

func NewTestAppForCalendar() *TestAppForCalendar {
	storageApp, _ := storage.InitStorage(storage.Config{
		Type: "memory",
	})

	return &TestAppForCalendar{
		Logger:  logger.New("debug"),
		Storage: storageApp,
	}
}

func (t *TestAppForCalendar) Setup() error {
	t.App = &app.App{
		Logger:  t.Logger,
		Storage: t.Storage,
	}

	handlers := internalhttp.NewCalendarHandlers(t.App, t.App.Logger)

	t.Server = internalhttp.NewServer(t.App, t.App.Logger, internalhttp.ServerConfig{
		Host:              "localhost",
		Port:              "8080",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}, handlers)

	go func() {
		_ = t.Server.Start(context.Background())
	}()

	return nil
}

func (t *TestAppForCalendar) Teardown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_ = t.Server.Stop(ctx)
}
