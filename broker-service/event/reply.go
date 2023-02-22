package event

import (
	"context"
	"log"
	"strings"
	"time"

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
		ShouldAutoAck: false,
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

	var forever chan struct{}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for rsp := range consumer {
			log.Printf(" [REPLIER] Received a message: %s", rsp.Body)
			body := strings.ToUpper(string(rsp.Body))

			msg := amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: rsp.CorrelationId,
				Body:          []byte(body),
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
			log.Printf(" [REPLIER] Sent response %s\n", body)

			rsp.Ack(false)
		}
	}()

	log.Printf(" [REPLIER] Awaiting requests")
	<-forever
}
