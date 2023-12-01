package main

import (
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/rs/zerolog/log"
)

func main() {
	util.Init_logger()
	log.Info().Msg("Starting Webserver...")
}
