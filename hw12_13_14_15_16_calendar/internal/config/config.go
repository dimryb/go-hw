package config

import (
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv" //nolint: depguard
)

type (
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	}

	Database struct {
		Type           string        `yaml:"type"`
		DSN            string        `yaml:"dsn"`
		MigrationsPath string        `yaml:"migrations" env:"MIGRATIONS_PATH"`
		Timeout        time.Duration `yaml:"timeout"`
	}
)

func Load(configPath string, target any) error {
	err := cleanenv.ReadConfig(path.Join("./", configPath), target)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(target)
	if err != nil {
		return fmt.Errorf("error updating env: %w", err)
	}

	return nil
}
