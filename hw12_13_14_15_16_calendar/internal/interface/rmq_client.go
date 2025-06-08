package interfaces

//go:generate mockgen -source=rmq_client.go -package=mocks -destination=../../mocks/mock_rmq_client.go
type RmqClient interface {
	Publish(routingKey string, body []byte) error
	Close() error
}
