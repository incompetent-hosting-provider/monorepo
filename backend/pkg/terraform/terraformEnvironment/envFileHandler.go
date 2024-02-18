package terraformenvironment

import (
	"incompetent-hosting-provider/backend/pkg/terraform/helper"
	"os"
	"reflect"

	"encoding/json"

	log "github.com/rs/zerolog/log"
)

func ReadRawEnv(path string, env_class reflect.Type) ([]byte, error) {
	// Read the tfenv file
	tf_env_raw, err := os.ReadFile(path)
	if err != nil {
		log.Error().Msgf("Error reading tfenv file: %s", err)
		return nil, err
	}

	// Unmarshal the tfenv file
	// Validate the tfenv file, by checking if the number of keys is the same as the number of keys in the CustomEnv struct
	var result map[string]interface{}
	json.Unmarshal([]byte(tf_env_raw), &result)

	refname := env_class.Name()
	refnum := env_class.NumField()
	print(refname, refnum)

	// Write a warning if there are more parameters than specified in the struct
	if len(result) > env_class.NumField() {
		log.Warn().Msgf("There are more parameters than expected in the tfenv file")
	}

	// Write a error if there are less parameters than specified in the struct
	if len(result) < env_class.NumField() {
		return nil, helper.NewCustomError("There are less parameters than expected in the tfenv file")
	}

	return tf_env_raw, nil
}

func WriteRawEnv(path string, tf_env any) error {
	// Marshal the tfenv file
	tf_env_raw, err := json.Marshal(tf_env)
	if err != nil {
		log.Error().Msgf("Error marshalling tfenv of type %s: %s", reflect.TypeOf(tf_env), err)
		return err
	}

	// Write the tfenv file
	err = os.WriteFile(path, tf_env_raw, 0644)
	if err != nil {
		log.Error().Msgf("Error writing tfenv file: %s", err)
		return err
	}

	return nil
}
