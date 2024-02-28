package mqhandler

import (
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/rs/zerolog/log"
)

func forwardQueueToChannel[T CustomContainerStartEvent | PresetContainerStartEvent | StopContainerEvent](ch *amqp.Channel, queueName string, targetChannel chan T) {
	q, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Error().Msgf("Could not declare queue due to an error: %v", err)
		// This is a go routine, no error is returned
		return
	}

	qc, err := ch.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		log.Error().Msgf("Could not consume queue due to an error: %v", err)
		return
	}

	// Serialize and pass into new channel
	for {
		in := <-qc
		out, err := parseReceivedEvent[T](in.Body)
		if err == nil {
			log.Debug().Msgf("Received event %v", out)
			targetChannel <- out
		}
	}
}
