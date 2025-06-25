package config

import "time"

type (
	SchedulerConfig struct {
		RabbitMQ  `yaml:"rabbitmq"`
		Database  `yaml:"database"`
		Scheduler `yaml:"scheduler"`
		Log       `yaml:"log"`
	}

	Scheduler struct {
		Interval        time.Duration `yaml:"interval" env:"INTERVAL"`
		RetentionPeriod time.Duration `yaml:"retentionPeriod"`
	}
)

func NewSchedulerConfig(configPath string) (*SchedulerConfig, error) {
	cfg := &SchedulerConfig{}
	err := Load(configPath, cfg)
	return cfg, err
}
