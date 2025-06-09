package config

type (
	SenderConfig struct {
		RabbitMQ  `yaml:"rabbitmq"`
		Log       `yaml:"log"`
		QueueName string `yaml:"queueName"`
	}
)
