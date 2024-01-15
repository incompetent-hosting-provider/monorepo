package auth

import (
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware(c *gin.Context) {
	log.Debug().Msg("Request passed auth middleware")

	// Note: This is a temporary implementatio and will be replaced with keycloak integration
	token := c.GetHeader("Authorization")

	if token == "" {
		util.ThrowUnauthorizedException(c, "Unauthorized")
		return
	}

	jwksURL := "http://localhost:8080/realms/ihp-realm/protocol/openid-connect/certs"

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.NewDefault([]string{jwksURL})

	if err != nil {
		log.Warn().Msg("Creating keyfunction not possible")
		util.ThrowUnauthorizedException(c, "Auth server not reachable")
	}

	parsed_token, err := jwt.Parse(token, jwks.Keyfunc)

	if err != nil {
		log.Warn().Msgf("Parsing failed: %v", err)
		util.ThrowUnauthorizedException(c, "Invalid token")
	}

	if claims, ok := parsed_token.Claims.(jwt.MapClaims); ok {
		c.Header("user-id", claims["sub"].(string))
	} else {
		log.Warn().Msg("Parsing claim failed")
		util.ThrowUnauthorizedException(c, "Invalid token")
	}

	c.Header("user-id", token)

	c.Next()
}
