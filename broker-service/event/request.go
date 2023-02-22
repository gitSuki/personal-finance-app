package event

import (
	"context"
	"log"
	"time"

	"github.com/gitsuki/finance/broker/util"
	"github.com/google/uuid"
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

type publishParams struct {
	Context      context.Context
	ExchangeName string // exchange (exchange to publish the message to)
	RoutingKey   string // routing key (routing key for the message. used by the exchange to determine which queues the message should be delivered to)
	IsMandatory  bool   // mandatory (whether the message must be routed to at least one queue. if true and the message cannot be routed, the message will be returned to the publisher)
	IsImmediate  bool   // immediate (whether the message should be delivered immediately. if set to true and the message cannot be delivered to a consumer immediately, it will be returned to the publisher)
	Message      amqp.Publishing
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
		IsExclusive:   false,
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

	corrId := uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := amqp.Publishing{
		ContentType:   "text/plain",
		CorrelationId: corrId,
		Body:          data,
	}

	publishArgs := &publishParams{
		Context:      ctx,
		ExchangeName: "",
		RoutingKey:   queue.Name,
		IsMandatory:  false,
		IsImmediate:  false,
		Message:      msg,
	}
	err = ch.PublishWithContext(
		publishArgs.Context,
		publishArgs.ExchangeName,
		publishArgs.RoutingKey,
		publishArgs.IsMandatory,
		publishArgs.IsImmediate,
		publishArgs.Message,
	)
	if err != nil {
		log.Panic("[panic] unable to publish message", err)
	}
	log.Printf(" [REQUESTER] Sent %s\n", message)

	for response := range consumer {
		if corrId == response.CorrelationId {
			body := string(response.Body)
			log.Printf(" [REQUESTER] Recieved response  %s", body)
		}
	}
}
