package interfaces

type RmqClient interface {
	Publish(routingKey string, body []byte) error
	Close() error
}
