package main

import (
	"log"

	"github.com/kchimev/locations-api/authentication/internal/constants"

	"github.com/streadway/amqp"
)

func setupRabbit() (*amqp.Connection, *amqp.Channel, error) {
	mqcon, err := amqp.Dial(constants.RabbitURL)
	if err != nil {
		return nil, nil, err
	}

	mqchan, err := mqcon.Channel()
	if err != nil {
		mqcon.Close()
		return nil, nil, err
	}

	err = mqchan.ExchangeDeclare(
		constants.RabbitExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		mqcon.Close()
		mqchan.Close()
		return nil, nil, err
	}

	queues := []struct {
		name       string
		routingKey string
	}{
		{constants.AuthenticateQueueName, constants.AuthenticateQueueKey},
		{constants.GenerateQueueName, constants.GenerateQueueKey},
	}

	for _, q := range queues {
		queue, err := mqchan.QueueDeclare(
			q.name,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			mqcon.Close()
			mqchan.Close()
			return nil, nil, err
		}

		err = mqchan.QueueBind(
			queue.Name,
			q.routingKey,
			constants.RabbitExchange,
			false,
			nil,
		)
		if err != nil {
			mqcon.Close()
			mqchan.Close()
			return nil, nil, err
		}
	}

	log.Println("RabbitMQ setup completed successfully")
	return mqcon, mqchan, nil
}

func (a *app) consumeQueue(queueName string, handler func(amqp.Delivery)) {
	msgs, err := a.mqchan.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages from %s: %v", queueName, err)
	}

	for msg := range msgs {
		handler(msg)
	}
}

func (a *app) handleAuthenticate(msg amqp.Delivery) {
	token := string(msg.Body)

	err := a.checkToken(token)
	response := "Valid"
	if err != nil {
		response = "Invalid"
	}

	if msg.ReplyTo != "" {
		a.mqchan.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(response),
			CorrelationId: msg.CorrelationId,
		})
	}
}

func (a *app) handleGenerate(msg amqp.Delivery) {
	response, err := a.generateToken()
	if err != nil {
		response = ""
	}

	if msg.ReplyTo != "" {
		a.mqchan.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(response),
			CorrelationId: msg.CorrelationId,
		})
	}
}
