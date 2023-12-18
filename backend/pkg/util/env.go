package util

import (
	"os"

	"github.com/rs/zerolog/log"
)

func GetStringEnvWithDefault(varName string, defaultValue string) string {
	val, exists := os.LookupEnv(varName)
	if !exists {
		log.Info().Msgf("Could not get env var %s defaulting to %s", varName, defaultValue)
		return defaultValue
	}
	return val
}
