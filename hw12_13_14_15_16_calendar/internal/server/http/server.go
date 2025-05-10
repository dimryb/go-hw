package internalhttp

import (
	"context"
)

type Server struct {
	logger Logger
	app    Application
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
