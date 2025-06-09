package calendar

import (
	"context"
	"fmt"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/config"
	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/http"
)

type Calendar struct {
	app  i.Application
	logg i.Logger
	cfg  *config.CalendarConfig
}

func NewCalendar(app i.Application, logger i.Logger, cfg *config.CalendarConfig) *Calendar {
	return &Calendar{
		app:  app,
		logg: logger,
		cfg:  cfg,
	}
}

func (s *Calendar) Run(ctx context.Context) error {
	handlers := internalhttp.NewCalendarHandlers(s.app, s.logg)

	server := internalhttp.NewServer(s.app, s.logg, internalhttp.ServerConfig{
		Host:              s.cfg.HTTP.Host,
		Port:              s.cfg.HTTP.Port,
		ReadTimeout:       s.cfg.HTTP.ReadTimeout,
		WriteTimeout:      s.cfg.HTTP.WriteTimeout,
		IdleTimeout:       s.cfg.HTTP.IdleTimeout,
		ReadHeaderTimeout: s.cfg.HTTP.ReadHeaderTimeout,
	}, handlers)

	if s.cfg.GRPC.Enable {
		go func() {
			s.logg.Debugf("gRPC server starting..")
			grpcServer := grpc.NewServer(
				s.app,
				grpc.ServerConfig{
					Port: s.cfg.GRPC.Port,
				},
				s.logg,
			)
			if err := grpcServer.Run(); err != nil {
				s.logg.Fatalf("Failed to start gRPC server: %s", err.Error())
			}
			s.logg.Debugf("gRPC server started")
		}()
	}

	go func() {
		<-ctx.Done()
		s.logg.Infof("Stopping HTTP server...")
		if err := server.Stop(context.Background()); err != nil {
			s.logg.Errorf("Failed to stop http server: %s", err.Error())
		}
	}()

	s.logg.Infof("calendar is running...")

	if err := server.Start(ctx); err != nil {
		return fmt.Errorf("failed to start http server: %s", err.Error())
	}
	return nil
}
