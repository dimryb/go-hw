package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/common"
	memorystorage "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type Config struct {
	Type           string
	DSN            string
	MigrationsPath string
	Timeout        time.Duration
	Migration      bool
}

func InitStorage(cfg Config) (storagecommon.EventStorage, error) {
	switch cfg.Type {
	case "memory":
		return memorystorage.New(), nil
	case "postgres":
		sqlStorage := sqlstorage.New(sqlstorage.Config{
			StorageType:    cfg.Type,
			DSN:            cfg.DSN,
			MigrationsPath: cfg.MigrationsPath,
		})

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		if err := sqlStorage.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		if cfg.Migration {
			if err := sqlStorage.Migrate(); err != nil {
				return nil, fmt.Errorf("failed to migrate database: %w", err)
			}
		}

		return sqlStorage, nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
	}
}
