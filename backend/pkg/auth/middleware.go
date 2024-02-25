package auth

import (
	"fmt"
	"incompetent-hosting-provider/backend/pkg/constants"
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware struct {
	JWKS keyfunc.Keyfunc
}

func isJWKSEndpointReachable(url string) bool {
	_, err := http.Get(url)
	if err != nil {
		print(err.Error())
		return false
	}
	return true
}

// the init function is called exactly once in golang
func GetAuthMiddleware() AuthMiddleware {
	JWKS_URL := util.GetStringEnvWithDefault("KEYCLOAK_CERT_ENDPOINT_URL", "http://localhost:8080/realms/ihp-realm/protocol/openid-connect/certs")

	if !isJWKSEndpointReachable(JWKS_URL) {
		log.Warn().Msg("JWKS endpoint not reachable")
		panic("JWKS endpoint not reachable")
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.NewDefault([]string{JWKS_URL})

	if err != nil || jwks == nil {
		log.Warn().Msg("Creating keyfunction not possible")
		// Panic this...because there is no way the backend can operate in this state
		panic("COULD NOT ESTABLISH INITIAL CONNECTION TO AUTH SERVER")
	}
	log.Info().Msg("JWKS initialized")

	return AuthMiddleware{JWKS: jwks}
}

func (a *AuthMiddleware) AuthFunc(c *gin.Context) {
	log.Debug().Msg("Request entered auth middleware")

	// Note: This is a temporary implementatio and will be replaced with keycloak integration
	token := c.GetHeader("Authorization")

	// Remove "Bearer" from Authorization Header in Format Bearer <jwt key>
	token = strings.TrimPrefix(token, "Bearer ")

	if token == "" {
		fmt.Print("no token")
		util.ThrowUnauthorizedException(c, "Unauthorized")
		return
	}

	log.Warn().Msgf("%v", a.JWKS)

	parsed_token, err := jwt.Parse(token, a.JWKS.Keyfunc)

	if err != nil {
		log.Warn().Msgf("Parsing failed: %v", err)
		util.ThrowUnauthorizedException(c, "Invalid token")
		return
	}

	if claims, ok := parsed_token.Claims.(jwt.MapClaims); ok {

		if claims["sub"] == nil || claims["email"] == nil {
			log.Warn().Msg("JWT without expected claims reported. This could mean that someone is tampering with JWTs")
			util.ThrowInternalServerErrorException(c, "JWT claim did not include expected claims")
		}

		sub := fmt.Sprintf("%v", claims["sub"])
		email := fmt.Sprintf("%v", claims["email"])

		c.Request.Header.Add(constants.USER_ID_HEADER, sub)
		c.Request.Header.Add(constants.USER_EMAIL_HEADER, email)
	} else {
		log.Warn().Msg("Parsing claim failed")
		util.ThrowUnauthorizedException(c, "Invalid token")
		return
	}

	log.Debug().Msg("Request passed auth middleware")

	c.Next()
}
