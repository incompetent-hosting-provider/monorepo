package incompetenthostingprovider

import (
	"encoding/json"
	"reflect"

	tfenv "goterra/pkg/terraformEnvironment"

	helper "goterra/pkg/helper"

	log "github.com/rs/zerolog/log"
)

type Env struct {
	// MySQL docker name prefix
	MySQLDockerNamePrefix string `json:"tfDocker_mysql_name_prefix"`
	// MySQL uids
	MySQLUIDs []string `json:"tfDocker_uids"`
}

func CompareCustomEnv(tf_env_current Env, tf_env_new Env) bool {
	// Compare the MySQLDockerNamePrefix
	if tf_env_current.MySQLDockerNamePrefix != tf_env_new.MySQLDockerNamePrefix {
		return false
	}

	// Compare the MySQLUIDs
	if !reflect.DeepEqual(tf_env_current.MySQLUIDs, tf_env_new.MySQLUIDs) {
		return false
	}

	return true
}

func ReadCustomEnv(path string) (Env, error) {
	tf_env_raw, err := tfenv.ReadRawEnv(path, reflect.TypeOf(Env{}))
	if err != nil {
		return Env{}, err
	}

	// Unmarshal the tfenv file
	var tf_env Env
	err = json.Unmarshal(tf_env_raw, &tf_env)
	if err != nil {
		log.Error().Msgf("Error unmarshalling tfenv file: %s", err)
		return Env{}, err
	}

	return tf_env, nil
}

func WriteCustomEnv(path string, tf_env Env) error {
	// Marshal the tfenv file
	err := tfenv.WriteRawEnv(path, tf_env)
	if err != nil {
		return err
	}

	return nil
}

func varInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (c *Env) AddMySqlContainer(uid string, port_external int, mysql_root_password string) error {
	if !varInSlice(uid, c.GetMySQLUIDs()) {
		log.Debug().Msgf("Adding MySQL container with UID %s", uid)

		c.SetMySQLUIDs(append(c.GetMySQLUIDs(), uid))
	} else {
		return helper.NewCustomError("UID already exists")
	}

	return nil
}

func (c *Env) RemoveMySqlContainer(index int) error {
	if index < 0 || index >= len(c.GetMySQLUIDs()) {
		return helper.NewCustomError("Index out of range")
	}

	log.Debug().Msgf("Removing MySQL container with UID %s", c.GetMySQLUIDs()[index])

	c.SetMySQLUIDs(append(c.GetMySQLUIDs()[:index], c.GetMySQLUIDs()[index+1:]...))

	return nil
}

// Getter
func (c *Env) GetMySQLDockerNamePrefix() string {
	return c.MySQLDockerNamePrefix
}

func (c *Env) GetMySQLUIDs() []string {
	return c.MySQLUIDs
}

// Setter
func (c *Env) SetMySQLDockerNamePrefix(name_prefix string) {
	c.MySQLDockerNamePrefix = name_prefix
}

func (c *Env) SetMySQLUIDs(uids []string) {
	c.MySQLUIDs = uids
}
