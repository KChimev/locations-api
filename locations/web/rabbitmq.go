package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/kchimev/locations-api/locations/internal/constants"
	"github.com/streadway/amqp"
)

type RabbitConnection struct {
	mqChan *amqp.Channel
	mqCon  *amqp.Connection
}

func setupRabbit() (*amqp.Connection, *amqp.Channel, error) {
	mqcon, err := amqp.Dial(constants.RabbitURL)
	if err != nil {
		return nil, nil, err
	}

	mqchan, err := mqcon.Channel()
	if err != nil {
		return nil, nil, err
	}

	return mqcon, mqchan, nil
}

func (c *RabbitConnection) validateToken(token string) (bool, error) {
	replyQueue, err := c.mqChan.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		return false, err
	}

	correlationID := uuid.New().String()

	err = c.mqChan.Publish(
		constants.AuthenticationExchange,
		constants.AuthenticateQueueKey,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(token),
			ReplyTo:       replyQueue.Name,
			CorrelationId: correlationID,
		},
	)
	if err != nil {
		return false, err
	}

	msgs, err := c.mqChan.Consume(
		replyQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return false, err
	}

	timeout := time.After(constants.AuthenticationTimeout * time.Second)
	for {
		select {
		case msg := <-msgs:
			if msg.CorrelationId == correlationID {
				return string(msg.Body) == "Valid", nil
			}
		case <-timeout:
			return false, nil
		}
	}
}
