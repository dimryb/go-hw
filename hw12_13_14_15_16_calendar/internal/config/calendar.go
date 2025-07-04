package config

import (
	"time"
)

type (
	CalendarConfig struct {
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

	GRPC struct {
		Enable bool   `yaml:"enable"`
		Port   string `yaml:"port" env:"GRPC_PORT"`
	}
)

func NewCalendarConfig(configPath string) (*CalendarConfig, error) {
	cfg := &CalendarConfig{}
	err := Load(configPath, cfg)
	return cfg, err
}
