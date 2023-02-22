package event

import (
	"log"

	"github.com/gitsuki/finance/broker/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type queueParams struct {
	Name             string     // name (name of queue that will be declared)
	IsDurable        bool       // durable (the queue will survive a broker restart)
	ShouldAutoDelete bool       // auto-delete (queue is deleted when last consumer unsubscribes)
	IsExclusive      bool       // exclusive (used by only one connection and the queue will be deleted when that connection closes)
	ShouldNotWait    bool       // nowait (client shouldn't wait for a response from the server before proceeding with the next operation.)
	Arguments        amqp.Table // arguments
}

type consumeParams struct {
	QueueName     string     // queue (name of queue to consume from.)
	Consumer      string     // consumer (string that identifies the consumer)
	ShouldAutoAck bool       // auto-ack (whether or not messages should be automatically acknowledged and removed from the queue)
	IsExclusive   bool       // exclusive (used by only one connection and the queue will be deleted when that connection closes)
	ShouldNoLocal bool       // no-local (whether or not messages published by this consumer should be ignored by itself)
	ShouldNotWait bool       // nowait (client shouldn't wait for a response from the server before proceeding with the next operation)
	Arguments     amqp.Table // arguments
}

func SendRequest(config util.Config, message string) {
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

	queueArgs := &queueParams{
		Name:             "testing_queue",
		IsDurable:        false,
		ShouldAutoDelete: false,
		IsExclusive:      true,
		ShouldNotWait:    false,
		Arguments:        nil,
	}
	queue, err := ch.QueueDeclare(
		queueArgs.Name,
		queueArgs.IsDurable,
		queueArgs.ShouldAutoDelete,
		queueArgs.IsExclusive,
		queueArgs.ShouldAutoDelete,
		queueArgs.Arguments,
	)
	if err != nil {
		log.Panic("[panic] unable to declare a rabbitmq queue", err)
	}

	consumerArgs := &consumeParams{
		QueueName:     "testing_queue",
		Consumer:      "testing_consumer",
		ShouldAutoAck: true,
		IsExclusive:   true,
		ShouldNoLocal: false,
		ShouldNotWait: false,
		Arguments:     nil,
	}
	consumer, err := ch.Consume(
		consumerArgs.QueueName,
		consumerArgs.Consumer,
		consumerArgs.ShouldAutoAck,
		consumerArgs.IsExclusive,
		consumerArgs.ShouldNoLocal,
		consumerArgs.ShouldNotWait,
		consumerArgs.Arguments,
	)
	if err != nil {
		log.Panic("[panic] unable to declare a rabbitmq consumer", err)
	}
}
