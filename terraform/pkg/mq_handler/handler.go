package mqhandler

import (
	"fmt"
	"goterra/pkg/helper"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MqHandler struct {
	mqConn                           *amqp.Connection
	CustomContainerStartEventChannel chan CustomContainerStartEvent
	PresetContainerStartEventChannel chan PresetContainerStartEvent
	StopContainerEventChannel        chan StopContainerEvent
	channel                          *amqp.Channel
}

func (m *MqHandler) Init() {
	conn, err := amqp.Dial(helper.GetStringEnvWithDefault("MQ_CONNECTION_STRING", "amqp://guest:guest@localhost:5672/"))

	if err != nil {
		helper.HandleFatalError(err, fmt.Sprintf("Could not start connection to rabbitmq due to an error: %v", err))
	}

	m.mqConn = conn

	ch, err := conn.Channel()

	if err != nil {
		helper.HandleFatalError(err, "Could not open channel using connection due to an error")
	}

	m.channel = ch

	helper.HandleFatalError(err, "Could not open channel due to an error")

	m.CustomContainerStartEventChannel = make(chan CustomContainerStartEvent)
	m.PresetContainerStartEventChannel = make(chan PresetContainerStartEvent)
	m.StopContainerEventChannel = make(chan StopContainerEvent)

	go forwardQueueToChannel[CustomContainerStartEvent](ch, "CustomContainerStartQueue", m.CustomContainerStartEventChannel)
	go forwardQueueToChannel[PresetContainerStartEvent](ch, "PresetContainerStartQueue", m.PresetContainerStartEventChannel)
	go forwardQueueToChannel[StopContainerEvent](ch, "StopContainerQueue", m.StopContainerEventChannel)

}
