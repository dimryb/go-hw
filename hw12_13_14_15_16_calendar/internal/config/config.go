package config

import (
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv" //nolint: depguard
)

type (
	Config struct {
		HTTP     `yaml:"http"`
		Log      `yaml:"log"`
		Database `yaml:"database"`
	}

	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

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

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
