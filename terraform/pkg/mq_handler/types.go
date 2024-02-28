package mqhandler

import (
	"encoding/json"
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

type DeystroyContainerEvent struct {
	ContainerUUID string
	UserId        string
}

type UpdateInstanceEvent struct {
	ContainerUUID string
	UserId        string
	NewStatus     string
}

func parseReceivedEvent[T CustomContainerStartEvent | PresetContainerStartEvent | DeystroyContainerEvent](input []byte) (T, error) {
	var output T
	err := json.Unmarshal(input, &output)
	return output, err
}
