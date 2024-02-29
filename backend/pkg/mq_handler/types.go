package mq_handler

import (
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PresetContainerStartEvent struct {
	ContainerUUID string
	UserId        string
	PresetId      int
	ContainerEnv  map[string]string
}

type CustomContainerStartEvent struct {
	ContainerUUID     string
	UserId            string
	ContainerImage    string
	ContainerImageTag string
	ContainerEnv      map[string]string
	ContainerPorts    []int
}

type DeleteContainerEvent struct {
	ContainerUUID string
	UserId        string
}

type mqWrapper struct {
	prestContainerEventQueue         amqp.Queue
	customContainerEventQueue        amqp.Queue
	stopContainerEventQueue          amqp.Queue
	mqConn                           *amqp.Connection
	mqChann                          *amqp.Channel
	CustomContainerStartEventChannel chan CustomContainerStartEvent
	PresetContainerStartEventChannel chan PresetContainerStartEvent
	DeleteContainerEventChannel      chan DeleteContainerEvent
	updateInstanceEventChannel       <-chan amqp091.Delivery
}

type UpdateInstanceEvent struct {
	ContainerUUID string
	UserId        string
	NewStatus     string
}

func serializeEvent[T CustomContainerStartEvent | PresetContainerStartEvent | DeleteContainerEvent](input T) ([]byte, error) {
	return json.Marshal(input)
}

func unserializeEvent(input []byte) (UpdateInstanceEvent, error) {
	var output UpdateInstanceEvent
	err := json.Unmarshal(input, &output)
	return output, err
}
