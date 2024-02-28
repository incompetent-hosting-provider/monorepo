package mqhandler

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (m *MqHandler) publishUpdateInstanceStatusEvent(event UpdateInstanceEvent) {
	eventMarshalled, err := json.Marshal(event)

	if err != nil {
		log.Error().Msgf("Could not publish event to bus due to an error: %v", err)
	}
	m.channel.PublishWithContext(
		context.Background(),
		"",
		"UpdateInstanceQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventMarshalled,
		},
	)
}
