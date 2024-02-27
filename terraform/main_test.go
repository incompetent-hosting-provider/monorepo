package main_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	helper "goterra/pkg/helper"
	ihp "goterra/pkg/incompetentHostingProvider"

	log "github.com/rs/zerolog/log"
)

const tf_version string = "1.7.1"

var cwd, _ = os.Getwd()
var tf_bin_dir string = filepath.Join(cwd, "bins")
var tf_cwd_dir string = filepath.Join(cwd, "TerraDocker/test")

// Note: terraform.tfvars.json is the default and does not need to be included
var tf_envs_names []string = []string{
	filepath.Join(tf_cwd_dir, "terraform.tfvars.json"),
	filepath.Join(tf_cwd_dir, "creds.tfvars.json"),
}

func TestRunMainApplyTerraform(t *testing.T) {
	ihpTfBin := ihp.NewTfBin(tf_bin_dir, tf_version, tf_cwd_dir, tf_envs_names)

	_, err := ihpTfBin.ApplyTerraform([]ihp.DockerMySQL{}, []ihp.DockerMySQL{})
	helper.HandleError(err, "Error applying terraform")
}

func TestRunMainReadWriteEnv(t *testing.T) {
	// Read the terraform environment
	tf_env_filepath := tf_envs_names[0]
	tf_env, err := ihp.ReadCustomEnv(tf_env_filepath)
	helper.HandleFatalError(err, "Error reading terraform environment")

	// Set the terraform environment
	err = ihp.WriteCustomEnv(tf_env_filepath, tf_env)
	helper.HandleError(err, "Error setting terraform environment")
}

func TestRunMainAddRemoveContainers(t *testing.T) {
	ihpTfBin := ihp.NewTfBin(tf_bin_dir, tf_version, tf_cwd_dir, tf_envs_names)

	uids := []int{}
	s := "test"
	for i := range s {
		uids = append(uids, int(s[i]))
	}

	passwords := []string{}
	for i := range uids {
		passwords = append(passwords, "password"+strconv.Itoa(i))
	}

	uidsToRemove := uids
	uidsToAdd := uids
	passwordsForUidsToAdd := passwords

	for i := range uidsToRemove {
		_, err := ihp.RemoveIhpMySqlContainer(ihpTfBin, uidsToRemove[i])
		helper.HandleError(err, "Error removing mysql container")
	}
	for i := range uidsToAdd {
		_, err := ihp.AddIhpMySqlContainer(ihpTfBin, uidsToAdd[i], passwordsForUidsToAdd[i])
		helper.HandleError(err, "Error adding mysql container")
	}

	state, err := ihpTfBin.GetState()
	helper.HandleError(err, "Error getting terraform state")

	current_num_mysql_containers := state.Values.Outputs["current_num_mysql_containers"].Value.(json.Number)
	log.Debug().Msgf("current_num_mysql_containers: %v", current_num_mysql_containers)
}
