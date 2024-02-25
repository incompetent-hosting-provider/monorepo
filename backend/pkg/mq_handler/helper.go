package mq_handler

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func getMqQueueFromChann(c *amqp.Channel, name string)(amqp.Queue){
	q, err := c.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil{
		panic(err)
	}

	return q
}

