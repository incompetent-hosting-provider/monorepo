package incompetenthostingprovider

import (
	"encoding/json"
	"goterra/pkg/helper"

	tfenv "goterra/pkg/terraformEnvironment"

	log "github.com/rs/zerolog/log"
)

type CredsEnv struct {
	// MySQL root passwords
	MySQLRootPasswords []string `json:"mysql_credentials"`
}

func ReadCredsEnv(path string) (CredsEnv, error) {
	tf_env_raw, err := tfenv.ReadRawEnv(path, helper.GetType(CredsEnv{}))
	if err != nil {
		return CredsEnv{}, err
	}

	// Unmarshal the tfenv file
	var tf_env CredsEnv
	err = json.Unmarshal(tf_env_raw, &tf_env)
	if err != nil {
		log.Error().Msgf("Error unmarshalling tfenv file: %s", err)
		return CredsEnv{}, err
	}

	return tf_env, nil
}

func WriteCredsEnv(path string, tf_env CredsEnv) error {
	err := tfenv.WriteRawEnv(path, tf_env)
	if err != nil {
		return err
	}

	return nil
}

func (i *CredsEnv) AddMysqlRootPassword(mysql_root_password string) {
	i.MySQLRootPasswords = append(i.MySQLRootPasswords, mysql_root_password)
}

func (i *CredsEnv) RemoveMysqlRootPassword(index int) error {
	// Check if the index is out of range
	if index < 0 || index >= len(i.MySQLRootPasswords) {
		return helper.NewCustomError("Index out of range")
	}

	i.MySQLRootPasswords = append(i.MySQLRootPasswords[:index], i.MySQLRootPasswords[index+1:]...)

	return nil
}

func (i *CredsEnv) RemoveMysqlRootPasswordFromIndex(index int) error {
	// Check if the index is out of range
	if index < 0 || index >= len(i.MySQLRootPasswords) {
		return helper.NewCustomError("Index out of range")
	}

	i.MySQLRootPasswords = append(i.MySQLRootPasswords[:index], i.MySQLRootPasswords[index+1:]...)

	return nil
}

// Getter
func (i *CredsEnv) GetMySQLRootPasswords() []string {
	return i.MySQLRootPasswords
}

// Setter
func (i *CredsEnv) SetMySQLRootPasswords(mysql_root_passwords []string) {
	i.MySQLRootPasswords = mysql_root_passwords
}
