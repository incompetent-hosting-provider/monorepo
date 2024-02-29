package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"goterra/pkg/helper"
	ihp "goterra/pkg/incompetentHostingProvider"
	logging "goterra/pkg/logging"
	mqhandler "goterra/pkg/mq_handler"

	log "github.com/rs/zerolog/log"
)

const tf_version string = "1.7.1"

var cwd, _ = os.Getwd()
var tf_bin_dir string = filepath.Join(cwd, "bins")
var tf_cwd_dir string = filepath.Join(cwd, "TerraDocker/dev")

// Note: terraform.tfvars.json is the default and does not need to be included
var tf_envs_names []string = []string{
	filepath.Join(tf_cwd_dir, "terraform.tfvars.json"),
	filepath.Join(tf_cwd_dir, "creds.tfvars.json"),
}

var ihp_presets = map[string]int{
	"mysql": 1,
}

func listenForEvents(mq *mqhandler.MqHandler, mq_preset_container_events *[]mqhandler.PresetContainerStartEvent) {
	for {
		select {
		case event := <-mq.CustomContainerStartEventChannel:
			// TODO: do something with this
			log.Debug().Msgf("Received custom start event: %v", event)
			log.Warn().Msg("CustomContainerStartEventChannel not implemented")
			updateEvent := mqhandler.UpdateInstanceEvent{UserId: event.UserId , ContainerUUID: event.ContainerUUID, NewStatus: "Running"}
			mq.PublishUpdateInstanceStatusEvent(updateEvent)
		case event := <-mq.PresetContainerStartEventChannel:
			event_sanitized := event
			event_sanitized.ContainerEnv = map[string]string{"<redacted>": "<redacted>"} // Redact sensitive data
			log.Debug().Msgf("Received preset start event: %v", event_sanitized)
			*mq_preset_container_events = append(*mq_preset_container_events, event)
		case event := <-mq.DestroyContainerEventChannel:
			// TODO: do something with this
			log.Debug().Msgf("Received delete event: %v", event)
			log.Warn().Msg("DestroyContainerEventChannel not implemented")
			updateEvent := mqhandler.UpdateInstanceEvent{UserId: event.UserId , ContainerUUID: event.ContainerUUID, NewStatus: "Deleted"}
			mq.PublishUpdateInstanceStatusEvent(updateEvent)
		}
	}
}

func main() {
	logging.InitLogger()
	log.Info().Msg("Starting GoTerra")

	ihpTfBin := ihp.NewTfBin(tf_bin_dir, tf_version, tf_cwd_dir, tf_envs_names)
	ihpTfBin.InitTerraform()

	mq_preset_container_events := []mqhandler.PresetContainerStartEvent{}

	mq := mqhandler.MqHandler{}
	mq.Init()

	go listenForEvents(&mq, &mq_preset_container_events)

	for {
		time.Sleep(5 * time.Second)

		if len(mq_preset_container_events) > 0 {

			uidsToAdd := []string{}
			passwordsForUidsToAdd := []string{}

			log.Debug().Msgf("Number of MQ preset container events: %v", len(mq_preset_container_events))
			for i := range mq_preset_container_events {
				if mq_preset_container_events[i].PresetId == ihp_presets["mysql"] {
					log.Debug().Msgf("Processing MySQL preset container from MQ: UserId=%v ContainerId=%v", mq_preset_container_events[i].UserId, mq_preset_container_events[i].ContainerUUID)
					conatiner_id := mq_preset_container_events[i].UserId + mq_preset_container_events[i].ContainerUUID
					log.Info().Msgf("Attemting to create container with UID: %v", conatiner_id)

					uidsToAdd = append(uidsToAdd, conatiner_id)
					passwordsForUidsToAdd = append(passwordsForUidsToAdd, mq_preset_container_events[i].ContainerEnv["MYSQL_ROOT_PASSWORD"])

					// Remove from list
					mq_preset_container_events = append(mq_preset_container_events[:i], mq_preset_container_events[i+1:]...)
				}
			}

			for i := range uidsToAdd {
				_, err := ihp.AddIhpMySqlContainer(ihpTfBin, uidsToAdd[i], passwordsForUidsToAdd[i])
				if err != nil {
					helper.HandleError(err, "Error adding mysql container")
				} else {
					log.Info().Msgf("Added MySQL container with UID %s", uidsToAdd[i])
					// Send update event
					// Get userid and containeruuid from uid by relying on the length of uuids (36)
					event := mqhandler.UpdateInstanceEvent{UserId: uidsToAdd[i][:36], ContainerUUID: uidsToAdd[i][36:], NewStatus: "Running"}
					log.Debug().Msgf("Sending update event: %v", event)
					mq.PublishUpdateInstanceStatusEvent(event)
				}
			}

			state, err := ihpTfBin.GetState()
			helper.HandleError(err, "Error getting terraform state")

			if state.Values != nil {
				current_num_mysql_containers := state.Values.Outputs["current_num_mysql_containers"].Value.(json.Number)
				log.Debug().Msgf("current_num_mysql_containers: %v", current_num_mysql_containers)
			}
		}
	}

	/*
		uidsToRemove := []int{}

		for i := range uidsToRemove {
			_, err := ihp.RemoveIhpMySqlContainer(ihpTfBin, uidsToRemove[i])
			helper.HandleError(err, "Error removing mysql container")
		}

		// _, err = ihpTfBin.ApplyTerraform([]ihp.DockerMySQL{}, []ihp.DockerMySQL{})
		// helper.HandleError(err, "Error applying terraform")
	*/

	// Exit main
	log.Info().Msg("Exiting GoTerra")
	os.Exit(0)
}
