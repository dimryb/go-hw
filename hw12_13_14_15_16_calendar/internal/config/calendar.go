package config

import (
	"time"
)

type (
	Config struct {
		HTTP     `yaml:"http"`
		Log      `yaml:"log"`
		Database `yaml:"database"`
		GRPC     `yaml:"grpc"`
	}

	HTTP struct {
		Host              string        `yaml:"host" env:"HTTP_HOST"`
		Port              string        `yaml:"port" env:"HTTP_PORT"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		IdleTimeout       time.Duration `yaml:"idleTimeout"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
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

	GRPC struct {
		Enable bool   `yaml:"enable"`
		Port   string `yaml:"port" env:"GRPC_PORT"`
	}
)

func NewCalendarConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	err := Load(configPath, cfg)
	return cfg, err
}
