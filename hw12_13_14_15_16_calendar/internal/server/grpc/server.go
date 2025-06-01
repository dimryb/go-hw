package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/grpc/interceptors"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/proto/calendar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Application interface {
	CreateEvent(context.Context, types.Event) error
	UpdateEvent(context.Context, types.Event) error
	DeleteEvent(context.Context, string) error
	GetEventByID(context.Context, string) (types.Event, error)
	ListEvents(context.Context) ([]types.Event, error)
	ListEventsByUser(context.Context, string) ([]types.Event, error)
	ListEventsByUserInRange(context.Context, string, time.Time, time.Time) ([]types.Event, error)
}

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type Server struct {
	app Application
	cfg ServerConfig
	log Logger
}

type ServerConfig struct {
	Port string
}

func NewServer(app Application, cfg ServerConfig, log Logger) *Server {
	return &Server{
		app: app,
		cfg: cfg,
		log: log,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryLoggerInterceptor(s.log)),
	)
	calendar.RegisterCalendarServiceServer(grpcServer, NewCalendarService(s.app))

	reflection.Register(grpcServer)

	s.log.Infof("Starting gRPC server, port %s", s.cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
