package auth

import (
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware(c *gin.Context) {
	log.Debug().Msg("Request passed auth middleware")

	// Note: This is a temporary implementatio and will be replaced with keycloak integration
	userId := c.GetHeader("Authorization")

	if userId == "" {
		util.ThrowUnauthorizedException(c, "Unauthorized")
		return
	}

	c.Header("user-id", userId)

	c.Next()
}
