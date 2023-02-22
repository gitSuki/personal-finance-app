package event

import (
	"log"

	"github.com/gitsuki/finance/broker/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RecieveRequest(config util.Config) {
	conn, err := amqp.Dial(config.RabbitMQ)
	if err != nil {
		log.Panic("[panic] unable to connect to rabbitmq", err)
	}
	defer conn.Close() // isn't called until the current function has returned

	ch, err := conn.Channel()
	if err != nil {
		log.Panic("[panic] unable to open a rabbitmq channel", err)
	}
	defer ch.Close() // isn't called until the current function has returned

	queueArgs := &queueParams{
		Name:             "testing_queue",
		IsDurable:        false,
		ShouldAutoDelete: false,
		IsExclusive:      false,
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
		QueueName:     queue.Name,
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

	var forever chan struct{}

	log.Printf(" [*] Awaiting requests.")
	<-forever
}
