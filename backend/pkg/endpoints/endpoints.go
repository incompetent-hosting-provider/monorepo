package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ConfigureEndpoints() *gin.Engine {
	log.Info().Msg("Setting up gin")
	ginEngine := gin.New()
	configureMiddleWares(ginEngine)
	configureGetEndpoints(ginEngine)
	return ginEngine
}

func configureMiddleWares(ginEngine *gin.Engine) {
	log.Info().Msg("Setting up middlewares")
	// Add recovery middleware to stop errors from oozing out
	ginEngine.Use(gin.Recovery())
}

func configureGetEndpoints(ginEngine *gin.Engine) {
	log.Info().Msg("Setting up GET endpoints")
	ginEngine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
}
