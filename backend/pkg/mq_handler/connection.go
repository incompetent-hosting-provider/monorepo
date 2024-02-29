package mq_handler

import (
	db_instances "incompetent-hosting-provider/backend/pkg/db/tables/instances"
	"incompetent-hosting-provider/backend/pkg/util"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

var mq mqWrapper

func (q *mqWrapper) runHandler() {
	defer q.mqConn.Close()

	for {
		select {
		case event := <-q.CustomContainerStartEventChannel:
			eventBody, _ := serializeEvent[CustomContainerStartEvent](event)
			err := q.mqChann.PublishWithContext(
				context.Background(),
				"",
				q.customContainerEventQueue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        eventBody,
				},
			)
			if err != nil {
				log.Warn().Msgf("Could not send message due to an error: %v", err)
			}
		case event := <-q.PresetContainerStartEventChannel:
			eventBody, _ := serializeEvent[PresetContainerStartEvent](event)
			err := q.mqChann.PublishWithContext(
				context.Background(),
				"",
				q.prestContainerEventQueue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        eventBody,
				},
			)
			if err != nil {
				log.Warn().Msgf("Could not send message due to an error: %v", err)
			}
		case event := <-q.DeleteContainerEventChannel:
			eventBody, _ := serializeEvent[DeleteContainerEvent](event)
			err := q.mqChann.PublishWithContext(
				context.Background(),
				"",
				q.stopContainerEventQueue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        eventBody,
				},
			)
			if err != nil {
				log.Warn().Msgf("Could not send message due to an error: %v", err)
			}
		case event := <-q.updateInstanceEventChannel:
			log.Debug().Msg("Received update event")
			parsedEvent, err := unserializeEvent(event.Body)
			if err != nil {
				log.Warn().Msgf("Could not parse update event due to an error: %v", err)
				continue
			}

			err = db_instances.UpdateInstanceStatus(parsedEvent.UserId, parsedEvent.ContainerUUID, parsedEvent.NewStatus)
			if err != nil {
				log.Warn().Msgf("Could not update instance")
			}
		}
	}
}

func initConnection() error {
	mqConnectionString := util.GetStringEnvWithDefault("MQ_CONN_STRING", "amqp://guest:guest@localhost:5672/")

	mqConn, err := amqp.Dial(mqConnectionString)

	if err != nil {
		return err
	}

	mqChann, err := mqConn.Channel()

	if err != nil {
		return err
	}

	q, err := mqChann.QueueDeclare(
		"UpdateInstanceQueue",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Error().Msgf("Could not declare queue due to an error: %v", err)
		return err
	}

	qc, err := mqChann.Consume(
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
		return err
	}

	mq = mqWrapper{
		CustomContainerStartEventChannel: make(chan CustomContainerStartEvent),
		PresetContainerStartEventChannel: make(chan PresetContainerStartEvent),
		DeleteContainerEventChannel:      make(chan DeleteContainerEvent),
		mqConn:                           mqConn,
		mqChann:                          mqChann,
		customContainerEventQueue:        getMqQueueFromChann(mqChann, "CustomContainerStartQueue"),
		prestContainerEventQueue:         getMqQueueFromChann(mqChann, "PresetContainerStartQueue"),
		stopContainerEventQueue:          getMqQueueFromChann(mqChann, "StopContainerQueue"),
		updateInstanceEventChannel:       qc,
	}

	go mq.runHandler()
	return nil
}

func isConnAvailable() bool {
	if mq.mqConn != nil {
		return true
	}
	// If conn is not ready -> try to set it up
	err := initConnection()
	return err == nil
}
