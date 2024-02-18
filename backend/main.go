package main

import (
	"fmt"
	"incompetent-hosting-provider/backend/pkg/db"
	"incompetent-hosting-provider/backend/pkg/endpoints"
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Swagger definitions

//	@securityDefinitions.apikey BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Auth token security

func main() {
	// Load env
	godotenv.Load()
	util.InitLogger()
	err := db.InitDbConn()
	if err != nil {
		log.Fatal().Msgf("%s", err)
		panic(err)
	}
	log.Info().Msg("Connected to DynamoDb")
	log.Info().Msg("Starting Webserver...")
	ginEngine := endpoints.ConfigureEndpoints()
	// Dont trust any proxies MIGHT not be what we need for deployment
	ginEngine.SetTrustedProxies(nil)
	// No serve to setup application but not actually start the webserver
	// Mostly used for terraform setup
	if util.GetStringEnvWithDefault("NO_SERVE", "") != "" {
		log.Warn().Msg("NO_SERVE was set. Exiting...")
		return
	}
	port := util.GetStringEnvWithDefault("PORT", "8081")
	log.Info().Msgf("Starting on port: %s", port)
	ginEngine.Run(fmt.Sprintf(":%s", port))
}
