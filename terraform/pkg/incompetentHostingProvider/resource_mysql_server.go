package incompetenthostingprovider

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"goterra/pkg/helper"

	tfexec "github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"

	log "github.com/rs/zerolog/log"
)

func GetMySQLContainers(tf *tfexec.Terraform) (map[int]DockerMySQL, error) {
	mysql_containers := make(map[int]DockerMySQL)

	state, err := tf.Show(context.TODO())
	if err != nil {
		log.Error().Msgf("Error showing terraform plan: %s", err)
		return nil, err
	}

	if state.Values != nil {
		for _, container := range state.Values.RootModule.Resources {

			if container.Type == "docker_container" {
				index, err := container.Index.(json.Number).Int64()
				if err != nil {
					return nil, err
				}

				uid, err := strconv.Atoi(strings.Split(container.AttributeValues["name"].(string), "-")[2])
				if err != nil {
					return nil, err
				}
				port_external, err := container.AttributeValues["ports"].([]interface{})[0].(map[string]interface{})["external"].(json.Number).Int64()
				if err != nil {
					return nil, err
				}
				mysql_root_password := strings.Split(container.AttributeValues["env"].([]interface{})[1].(string), "=")[1]

				mysql_container := DockerMySQL{int(index), uid, int(port_external), mysql_root_password}

				mysql_containers[uid] = mysql_container
			}
		}
	}

	return mysql_containers, nil
}

func AddIhpMySqlContainer(tfbin *TerraformBinary, uid int, mysql_root_password string) (*tfjson.State, error) {
	log.Debug().Msgf("Adding container with UID %d", uid)
	// Checking the length of the UID is obligatory, the length of the datatype int is lower than 114
	// The maxium length of a docker container name is 128, 114 is the length minus the prefix
	if uid < 0 && len(strconv.Itoa(uid)) < 114 {
		log.Error().Msgf("UID %d is invalid", uid)
		return nil, helper.NewCustomError("Invalid UID")
	}

	tf, err := tfbin.GetInstance()
	if err != nil {
		return nil, err
	}

	// Check if the uid already exists, if no then add it to the create list
	if exists, err := containerWithUidExists(tf, uid); !exists && err == nil {
		containers_to_create := []DockerMySQL{
			{
				uid:                 uid,
				mysql_root_password: mysql_root_password,
			},
		}

		state, err := tfbin.ApplyTerraform(containers_to_create, []DockerMySQL{})
		if err != nil {
			return nil, err
		}

		return state, nil

	} else if err != nil {
		log.Error().Msgf("Error checking if container with UID %d exists: %s", uid, err)
		return nil, err

	} else {
		log.Error().Msgf("Container with UID %d already exists", uid)
		return nil, helper.NewCustomError("Container with UID already exists")

	}
}

func RemoveIhpMySqlContainer(tfbin *TerraformBinary, uid int) (*tfjson.State, error) {
	log.Debug().Msgf("Removing container with UID %d", uid)
	// Checking the length of the UID is obligatory, the length of the datatype int is lower than 114
	if uid < 0 && len(strconv.Itoa(uid)) < 114 {
		log.Error().Msgf("UID %d is invalid", uid)
		return nil, helper.NewCustomError("Invalid UID")
	}

	tf, err := tfbin.GetInstance()
	if err != nil {
		return nil, err
	}

	// Check if the container exists, if yes then add it to the destroy list
	if exists, err := containerWithUidExists(tf, uid); exists && err == nil {
		container, err := GetContainerWithUid(tf, uid)
		if err != nil {
			return nil, err
		}

		containers_to_destroy := []DockerMySQL{
			container,
		}

		state, err := tfbin.ApplyTerraform([]DockerMySQL{}, containers_to_destroy)
		if err != nil {
			return nil, err
		}

		return state, nil

	} else if err != nil {
		log.Error().Msgf("Error checking if container with UID %d exists: %s", uid, err)
		return nil, err

	} else {
		log.Error().Msgf("Container with UID %d does not exist", uid)
		return nil, helper.NewCustomError("Container with UID does not exist")

	}
}

func GetContainerWithUid(tf *tfexec.Terraform, uid int) (DockerMySQL, error) {
	containers, err := GetMySQLContainers(tf)
	if err != nil {
		return DockerMySQL{}, err
	}

	if exists, err := containerWithUidExists(tf, uid); exists {
		return containers[uid], err
	} else {
		return DockerMySQL{}, helper.NewCustomError("Container with UID does not exist")
	}
}

func containerWithUidExists(tf *tfexec.Terraform, uid int) (bool, error) {
	containers, err := GetMySQLContainers(tf)
	if err != nil {
		return false, err
	}

	if _, exists := containers[uid]; exists {
		return true, nil
	} else {
		return false, nil
	}
}
