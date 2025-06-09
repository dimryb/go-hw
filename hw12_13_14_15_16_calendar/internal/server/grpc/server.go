package grpc

import (
	"fmt"
	"net"

	i "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/interface"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/grpc/interceptors"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/proto/calendar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	app i.Application
	cfg ServerConfig
	log i.Logger
}

type ServerConfig struct {
	Port string
}

func NewServer(app i.Application, cfg ServerConfig, log i.Logger) *Server {
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
