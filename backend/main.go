package main

import (
	"fmt"
	"incompetent-hosting-provider/backend/pkg/db"
	"incompetent-hosting-provider/backend/pkg/endpoints"
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	util.InitLogger()
}

// Swagger definitions

// @securityDefinitions.apikey BearerAuth
// @in							header
// @name						Authorization
// @description				Auth token security
func main() {
	err := db.InitDbConn()
	if err != nil {
		log.Fatal().Msgf("%s", err)
		panic(err)
	}
	log.Info().Msg("Connected to DynamoDb")
	log.Info().Msg("Starting Webserver...")
	ginEngine := endpoints.ConfigureEndpoints()
	// Dont trust any proxies MIGHT not be what we need for deployment
	err = ginEngine.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal().Msgf("Could not add trusted proxies due to an error: %v", err)
	}
	port := util.GetStringEnvWithDefault("PORT", "8081")
	log.Info().Msgf("Starting on port: %s", port)
	err = ginEngine.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
