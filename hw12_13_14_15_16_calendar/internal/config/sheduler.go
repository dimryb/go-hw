package config

import "time"

type (
	SchedulerConfig struct {
		RabbitMQ          `yaml:"rabbitmq"`
		SchedulerDatabase `yaml:"database"`
		Scheduler         `yaml:"scheduler"`
	}

	RabbitMQ struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Exchange string `yaml:"exchange"`
	}

	SchedulerDatabase struct {
		Type    string        `yaml:"type"`
		DSN     string        `yaml:"dsn"`
		Timeout time.Duration `yaml:"timeout"`
	}

	Scheduler struct {
		Interval        time.Duration `yaml:"interval"`
		RetentionPeriod time.Duration `yaml:"retention_period"`
	}
)

func NewSchedulerConfig(configPath string) (*SchedulerConfig, error) {
	cfg := &SchedulerConfig{}
	err := Load(configPath, cfg)
	return cfg, err
}
