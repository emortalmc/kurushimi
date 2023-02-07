package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"kurushimi/internal/config"
)

const rabbitMqUriFormat = "amqp://%s:%s@%s:5672"

func NewConnection(cfg config.RabbitMQConfig) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf(rabbitMqUriFormat, cfg.Username, cfg.Password, cfg.Host))
}
