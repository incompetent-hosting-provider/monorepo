package util

import (
	"flag"
	"io/fs"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func IsTestRun() bool {
	if flag.Lookup("test.v") != nil || strings.HasSuffix(os.Args[0], ".test") {
		return true
	}
	return false
}

func CreateDirFromAbolutePathIfNotExist(path string, permission int) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, fs.FileMode(permission))
		if err != nil {
			log.Warn().Msgf("Could not create dir %s due to an error: %v", path, err)
			return
		}
		log.Info().Msgf("Created %s", path)
		return
	}
}
