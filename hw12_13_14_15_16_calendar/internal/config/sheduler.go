package config

import "time"

type (
	SchedulerConfig struct {
		RabbitMQ  `yaml:"rabbitmq"`
		Database  `yaml:"database"`
		Scheduler `yaml:"scheduler"`
		Log       `yaml:"log"`
	}

	RabbitMQ struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Exchange string `yaml:"exchange"`
	}

	Scheduler struct {
		Interval        time.Duration `yaml:"interval"`
		RetentionPeriod time.Duration `yaml:"retentionPeriod"`
	}
)

func NewSchedulerConfig(configPath string) (*SchedulerConfig, error) {
	cfg := &SchedulerConfig{}
	err := Load(configPath, cfg)
	return cfg, err
}
