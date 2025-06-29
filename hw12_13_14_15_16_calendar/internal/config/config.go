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
		DSN            string        `yaml:"dsn" env:"DATABASE_DSN"`
		MigrationsPath string        `yaml:"migrations" env:"MIGRATIONS_PATH"`
		Migrate        bool          `yaml:"migrate" env:"MIGRATE"`
		Timeout        time.Duration `yaml:"timeout"`
	}

	RabbitMQ struct {
		Host     string `yaml:"host" env:"RABBIT_HOST"`
		Port     string `yaml:"port" env:"RABBIT_PORT"`
		User     string `yaml:"user" env:"RABBIT_USER"`
		Password string `yaml:"password" env:"RABBIT_PASSWORD"`
		Exchange string `yaml:"exchange" env:"RABBIT_EXCHANGE"`
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
