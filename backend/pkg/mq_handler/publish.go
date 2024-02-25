package mq_handler

import "errors"

func PublishPresetContainerStartEvent(e PresetContainerStartEvent) error {
	if !isConnAvailable() {
		return errors.New("could not publish event at this time")
	}
	mq.PresetContainerStartEventChannel <- e
	return nil
}

func PublishCustomContainerStartEvent(e CustomContainerStartEvent) error {
	if !isConnAvailable() {
		return errors.New("could not publish event at this time")
	}

	mq.CustomContainerStartEventChannel <- e
	return nil
}

func PublishDeleteContainerEvent(e DeleteContainerEvent) error {
	if !isConnAvailable() {
		return errors.New("could not publish event at this time")
	}
	mq.DeleteContainerEventChannel <- e
	return nil
}
