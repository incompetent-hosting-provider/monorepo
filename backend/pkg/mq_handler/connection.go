package mq_handler

import (
	"incompetent-hosting-provider/backend/pkg/util"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/context"
)

var mq mqWrapper

func (q *mqWrapper) runHandler(){
	defer q.mqConn.Close()

	for {

		select{
		case event := <-q.CustomContainerStartEventChannel:
			eventBody, _ := serializeEvent[CustomContainerStartEvent](event)
			q.mqChann.PublishWithContext(
				context.Background(),
				"",
				q.customContainerEventQueue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body: eventBody,
				},
			)
		case event := <- q.PresetContainerStartEventChannel:
			eventBody, _ := serializeEvent[PresetContainerStartEvent](event)
			q.mqChann.PublishWithContext(
				context.Background(),
				"",
				q.prestContainerEventQueue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body: eventBody,
				},
			)
		case event := <- q.DeleteContainerEventChannel:
						eventBody, _ := serializeEvent[DeleteContainerEvent](event)
			q.mqChann.PublishWithContext(
				context.Background(),
				"",
				q.stopContainerEventQueue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body: eventBody,
				},
			)
		}
	}
}


func initConnection() error {
	mqConnectionString := util.GetStringEnvWithDefault("MQ_CONN_STRING", "amqp://guest:guest@localhost:5672/")

	mqConn, err := amqp.Dial(mqConnectionString)

	if err != nil{
		return err
	}

	mqChann, err := mqConn.Channel()

	if err != nil{
		return err
	}

	mq = mqWrapper{
		CustomContainerStartEventChannel: make(chan CustomContainerStartEvent),
		PresetContainerStartEventChannel: make(chan PresetContainerStartEvent),
		DeleteContainerEventChannel: make(chan DeleteContainerEvent),
		mqConn: mqConn,
		mqChann: mqChann,
		customContainerEventQueue: getMqQueueFromChann(mqChann, "CustomContainerStartQueue"),
		prestContainerEventQueue: getMqQueueFromChann(mqChann, "PresetContainerStartQueue"),
		stopContainerEventQueue: getMqQueueFromChann(mqChann, "StopContainerQueue"),	
	}

	go mq.runHandler()
	return nil
}

func isConnAvailable() bool{
	if mq.mqConn != nil{
		return true
	}
	// If conn is not ready -> try to set it up
	err := initConnection()
	if err != nil{
		return false
	}
	return true

}