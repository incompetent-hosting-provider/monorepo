package terraformenvironment

import (
	"encoding/json"
	"reflect"

	"github.com/rs/zerolog/log"
)

type EnvInterface interface {
	ReadEnv(path string) ([]byte, error)
	WriteEnv(path string, tf_env any) error
}

type Env struct {
	path string `json:"path"`
}

func NewEnv(path string) *Env {
	return &Env{
		path: path,
	}
}

func (e *Env) ReadEnv() (Env, error) {
	tf_env_raw, err := ReadRawEnv(e.path, reflect.TypeOf(Env{}))
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

func (e *Env) WriteEnv(tf_env Env) error {
	err := WriteRawEnv(e.path, tf_env)
	if err != nil {
		return err
	}

	return nil
}
