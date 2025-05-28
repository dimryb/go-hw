package grpc

import (
	"fmt"
	"net"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/proto/calendar"
	"google.golang.org/grpc"
)

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type Server struct {
	cfg ServerConfig
	log Logger
}

type ServerConfig struct {
	Port string
}

func NewServer(cfg ServerConfig, log Logger) *Server {
	return &Server{cfg: cfg, log: log}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	calendar.RegisterCalendarServiceServer(grpcServer, NewCalendarService())

	s.log.Infof("Starting gRPC server, port %s", s.cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
