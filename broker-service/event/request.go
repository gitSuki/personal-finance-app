package event

import (
	"log"

	"github.com/gitsuki/finance/broker/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendRequest(config util.Config, message string) error {
	data := []byte(message)

	conn, err := amqp.Dial(config.RabbitMQ)
	if err != nil {
		log.Panic("[panic] unable to connect to rabbitmq", err)
	}
	defer conn.Close() // isn't called until the current function has returned

	ch, err := conn.Channel()
	if err != nil {
		log.Panic("[panic] unable to open a rabbitmq channel", err)
	}
	defer ch.Close()
}
