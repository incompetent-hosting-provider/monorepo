package auth

import (
	"errors"
	"fmt"
	"incompetent-hosting-provider/backend/pkg/constants"
	"incompetent-hosting-provider/backend/pkg/util"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware struct {
	JWKS keyfunc.Keyfunc
}

// the init function is called exactly once in golang
func GetAuthMiddleware() AuthMiddleware {
	JWKS_URL := util.GetStringEnvWithDefault("KEYCLOAK_CERT_ENDPOINT_URL", "http://localhost:8080/realms/ihp-realm/protocol/openid-connect/certs")
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

func getKeyFromHeader(claim jwt.MapClaims, key string) (string, error) {
	val, ok := claim["foo"]
	// If the key exists
	if ok {
		return fmt.Sprint(val), nil
	}
	return "", errors.New("JWT claim did not contain expected key")
}

func (a *AuthMiddleware) AuthFunc(c *gin.Context) {
	log.Debug().Msg("Request passed auth middleware")

	// Note: This is a temporary implementatio and will be replaced with keycloak integration
	token := c.GetHeader("Authorization")

	// Remove "Bearer" from Authorization Header in Format Bearer <jwt key>
	token = strings.TrimPrefix(token, "Bearer ")

	if token == "" {
		fmt.Print("no token")
		util.ThrowUnauthorizedException(c, "Unauthorized")
		return
	}

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
	}

	c.Next()
}
