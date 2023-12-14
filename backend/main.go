package main

import (
	"fmt"
	"incompetent-hosting-provider/backend/pkg/endpoints"
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load env
	godotenv.Load()
	util.InitLogger()
	log.Info().Msg("Starting Webserver...")
	ginEngine := endpoints.ConfigureEndpoints()
	// Dont trust any proxies MIGHT not be what we need for deployment
	ginEngine.SetTrustedProxies(nil)
	port := util.GetStringEnvWithDefault("PORT", "8081")
	log.Info().Msgf("Starting on port: %s", port)
	ginEngine.Run(fmt.Sprintf(":%s", port))
}
