package mq_handler

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PresetContainerStartEvent struct{
	ContainerUUID string
	UserId string
	PresetId string
}

type CustomContainerStartEvent struct{
	ContainerUUID string
	UserId string
	ContainerImage string
	ContainerImageTag string
	ContainerEnv map[string]string
	ContainerPorts []string
}

type DeleteContainerEvent struct{
	ContainerUUID string
	UserId string
}


type mqWrapper struct{
	prestContainerEventQueue amqp.Queue
	customContainerEventQueue amqp.Queue
	stopContainerEventQueue amqp.Queue
	mqConn *amqp.Connection
	mqChann *amqp.Channel
	CustomContainerStartEventChannel chan CustomContainerStartEvent 
	PresetContainerStartEventChannel chan PresetContainerStartEvent 
	DeleteContainerEventChannel chan DeleteContainerEvent
}


func serializeEvent[T CustomContainerStartEvent | PresetContainerStartEvent | DeleteContainerEvent](input T) ([]byte, error){
	return json.Marshal(input)
}