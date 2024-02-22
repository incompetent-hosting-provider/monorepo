package endpoints

import (
	"incompetent-hosting-provider/backend/pkg/auth"
	"incompetent-hosting-provider/backend/pkg/payment"
	"incompetent-hosting-provider/backend/pkg/user"

	docs "incompetent-hosting-provider/backend/docs"
	"net/http"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ConfigureEndpoints() *gin.Engine {
	log.Info().Msg("Setting up gin")
	ginEngine := gin.New()
	authMiddleware := auth.GetAuthMiddleware()
	configureGlobalMiddleWares(ginEngine)
	configureGetEndpoints(ginEngine, authMiddleware)
	configurePostEndpoints(ginEngine, authMiddleware)
	configureSwagger(ginEngine)
	return ginEngine
}

func configureGlobalMiddleWares(ginEngine *gin.Engine) {
	log.Info().Msg("Setting up middlewares")
	// Add recovery middleware to stop errors from oozing out
	ginEngine.Use(gin.Recovery())

	p := ginprometheus.NewPrometheus("gin")
	p.Use(ginEngine)
	p.SetListenAddress("/metrics")
}

func configureGetEndpoints(ginEngine *gin.Engine, authMiddleware auth.AuthMiddleware) {
	log.Info().Msg("Setting up GET endpoints")

	ginEngine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	ginEngine.GET("/payment", authMiddleware.AuthFunc, payment.CreditFetchHandler)
	ginEngine.GET("/user", authMiddleware.AuthFunc, user.UserFetchHandler)
}

func configurePostEndpoints(ginEngine *gin.Engine, authMiddleware auth.AuthMiddleware) {

	log.Info().Msg("Setting up POST endpoints")
	ginEngine.POST("/payment", authMiddleware.AuthFunc, payment.ChangeCreditHandler)
}

func configureSwagger(ginEngine *gin.Engine) {
	if gin.Mode() != gin.ReleaseMode {
		log.Info().Msgf("Serving with Swagger version %s", docs.SwaggerInfo.Version)

		ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	} else {
		log.Info().Msg("Not serving swagger in release mode")
	}
}
