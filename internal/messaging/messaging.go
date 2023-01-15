package messaging

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"kurushimi/internal/config/dynamic"
)

const amqpUrl = "amqp://%s:%s@%s:5672"

type Messenger struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

func NewRabbitMQ(config dynamic.Config) (*Messenger, error) {
	rConfig := config.RabbitMq
	conn, err := amqp091.Dial(fmt.Sprintf(amqpUrl, rConfig.Username, rConfig.Password, rConfig.Host))
	if err != nil {
		return nil, err
	}
	chann, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Messenger{
		Connection: conn,
		Channel:    chann,
	}, nil
}
