package sqlstorage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

type Config struct {
	StorageType    string
	DSN            string
	MigrationsPath string
}

type Storage struct {
	storageType    string
	dsn            string
	migrationsPath string
	db             *sqlx.DB
}

func New(cfg Config) *Storage {
	return &Storage{
		storageType: cfg.StorageType,
		dsn:         cfg.DSN,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Connect(s.storageType, s.dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", s.storageType, err)
	}

	s.db = db

	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping %s: %w", s.storageType, err)
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("failed to close DB connection: %w", err)
		}
		s.db = nil
	}
	return nil
}

func (s *Storage) Migrate() error {
	if err := goose.SetDialect(s.storageType); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(nil, s.migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
