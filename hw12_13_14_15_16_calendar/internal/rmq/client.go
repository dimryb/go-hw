package rmq

import (
	"github.com/streadway/amqp" //nolint:depguard
)

type Client interface {
	Publish(routingKey string, body []byte) error
	Consume(queueName string) (<-chan []byte, error)
	Close() error
}

type client struct {
	channel  *amqp.Channel
	exchange string
}

func NewClient(amqpURL, exchange string) (Client, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &client{
		channel:  ch,
		exchange: exchange,
	}, nil
}

func (c *client) Publish(routingKey string, body []byte) error {
	return c.channel.Publish(
		c.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *client) Consume(queueName string) (<-chan []byte, error) {
	queue, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := c.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	out := make(chan []byte)
	go func() {
		for msg := range msgs {
			out <- msg.Body
		}
	}()
	return out, nil
}

func (c *client) Close() error {
	return c.channel.Close()
}
