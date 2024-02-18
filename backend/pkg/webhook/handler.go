package webhook

import (
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type KeycloakWebookPayload struct {
	// This is either DELETE_ACCOUNT and REGISTER
	KeycloakEvent string `json:"keycloakEvent"`
	UserId        string `json:"userId"`
}

// godoc
// @Summary 					SPI webhook handler
//
// @Schemes
// @Description 				Handle keycloak SPI webhook events
// @Tags 						webhook
//
// @Success 					200 {object} payment.BalanceResponse
//
// @Failure						500 {object} util.ErrorResponse
//
// @Router /spi-webhook [get]
func WebhookHandler(c *gin.Context) {
	var receivedEvent KeycloakWebookPayload

	err := c.ShouldBindJSON(&receivedEvent)
	if err != nil {
		log.Info().Msg("Could not parse request body")
		util.ThrowBadRequestException(c, "Could not parse body.")
		return
	}
	log.Info().Msgf("Received event %s", receivedEvent.KeycloakEvent)
}
