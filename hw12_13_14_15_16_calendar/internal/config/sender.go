package config

type (
	SenderConfig struct {
		RabbitMQ  `yaml:"rabbitmq"`
		Log       `yaml:"log"`
		QueueName string `yaml:"queueName"`
	}
)

func NewSenderConfig(path string) (*SenderConfig, error) {
	cfg := &SenderConfig{}
	err := Load(path, cfg)
	return cfg, err
}
