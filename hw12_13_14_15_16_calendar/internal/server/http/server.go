package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	// Импортируем сгенерированный пакет docs для регистрации Swagger.
	_ "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/server/http/docs"
	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
	httpSwagger "github.com/swaggo/http-swagger" //nolint: depguard
)

type Server struct {
	app    Application
	logger Logger
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
	CreateEvent(context.Context, types.Event) (string, error)
	UpdateEvent(context.Context, types.Event) error
	DeleteEvent(context.Context, string) error
	GetEventByID(context.Context, string) (types.Event, error)
	ListEvents(context.Context) ([]types.Event, error)
	ListEventsByUser(context.Context, string) ([]types.Event, error)
	ListEventsByUserInRange(context.Context, string, time.Time, time.Time) ([]types.Event, error)
}

type ServerConfig struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

func NewServer(app Application, logger Logger, cfg ServerConfig, handlers *CalendarHandlers) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/event/create", handlers.CreateEvent)
	mux.HandleFunc("/event/update", handlers.UpdateEvent)
	mux.HandleFunc("/event/delete", handlers.DeleteEvent)
	mux.HandleFunc("/event/get", handlers.GetEventByID)
	mux.HandleFunc("/events/list", handlers.ListEvents)
	mux.HandleFunc("/events/user", handlers.ListEventsByUser)
	mux.HandleFunc("/events/range", handlers.ListEventsByUserInRange)

	mux.HandleFunc("/", handlers.helloHandler)

	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		httpSwagger.Handler()(w, r)
	})

	return &Server{
		logger: logger,
		app:    app,
		server: &http.Server{
			Handler:           loggingMiddleware(handlers.logger)(mux),
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
		cfg: cfg,
	}
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

func (s *Server) Handler() http.Handler {
	return s.server.Handler
}
