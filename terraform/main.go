package main

import (
	"encoding/json"
	"os"

	helper "goterra/pkg/helper"
	ihp "goterra/pkg/incompetentHostingProvider"
	logging "goterra/pkg/logging"

	log "github.com/rs/zerolog/log"
)

var cwd, _ = os.Getwd()
var tf_bin_dir string = cwd + "/bins"
var tf_version string = "1.7.1"
var tf_cwd_dir string = cwd + "/TerraDocker/dev"

// Note: terraform.tfvars.json is the default and does not need to be included
var tf_envs_names []string = []string{
	tf_cwd_dir + "/terraform.tfvars.json",
	tf_cwd_dir + "/creds.tfvars.json",
}

func main() {
	logging.InitLogger()
	log.Info().Msg("Starting GoTerra")

	ihpTfBin := ihp.NewTfBin(tf_bin_dir, tf_version, tf_cwd_dir, tf_envs_names)

	uidsToRemove := []int{}

	uidsToAdd := []int{}
	passwordsForUidsToAdd := []string{"password4"}

	for i := range uidsToRemove {
		_, err := ihp.RemoveIhpMySqlContainer(ihpTfBin, uidsToRemove[i])
		helper.HandleError(err, "Error removing mysql container")
	}

	for i := range uidsToAdd {
		_, err := ihp.AddIhpMySqlContainer(ihpTfBin, uidsToAdd[i], passwordsForUidsToAdd[i])
		helper.HandleError(err, "Error adding mysql container")
	}

	// _, err = ihpTfBin.ApplyTerraform([]ihp.DockerMySQL{}, []ihp.DockerMySQL{})
	// helper.HandleError(err, "Error applying terraform")

	state, err := ihpTfBin.GetState()
	helper.HandleError(err, "Error getting terraform state")

	if state.Values != nil {
		current_num_mysql_containers := state.Values.Outputs["current_num_mysql_containers"].Value.(json.Number)
		log.Debug().Msgf("current_num_mysql_containers: %v", current_num_mysql_containers)
	}

	// Exit main
	log.Info().Msg("Exiting GoTerra")
	os.Exit(0)
}
