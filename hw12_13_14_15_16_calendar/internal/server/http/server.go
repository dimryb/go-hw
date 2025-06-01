package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
)

type Server struct {
	logger Logger
	app    Application
	server *http.Server
	cfg    ServerConfig
}

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type Application interface {
	CreateEvent(context.Context, types.Event) error
	// TODO
}

type ServerConfig struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

func NewServer(logger Logger, app Application, cfg ServerConfig) *Server {
	mux := http.NewServeMux()
	srv := &Server{
		logger: logger,
		app:    app,
		server: &http.Server{
			Handler:           mux,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
		cfg: cfg,
	}

	mux.Handle("/", loggingMiddleware(logger)(http.HandlerFunc(srv.helloHandler)))
	return srv
}

func (s *Server) Start(_ context.Context) error {
	addr := net.JoinHostPort(s.cfg.Host, s.cfg.Port)
	s.server.Addr = addr

	s.logger.Infof(fmt.Sprintf("Starting HTTP server on %s", addr))
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Errorf("Failed to start HTTP server: " + err.Error())
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Infof("Stopping HTTP server")
	return s.server.Shutdown(ctx)
}

func (s *Server) helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello, world!"))
}
