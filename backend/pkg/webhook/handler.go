package webhook

import (
	"incompetent-hosting-provider/backend/pkg/util"
	db_payment "incompetent-hosting-provider/backend/pkg/db/tables"
	"net/http"

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

	// Log the received event type
	log.Info().Msgf("Received event %s", receivedEvent.KeycloakEvent)

	// React based on the event type
	switch receivedEvent.KeycloakEvent {
	case "REGISTER":
		handleUserRegistration(receivedEvent.UserId)
	case "DELETE_ACCOUNT":
		handleUserDeletion(receivedEvent.UserId)
	default:
		log.Warn().Msgf("Unsupported Keycloak event type: %s", receivedEvent.KeycloakEvent)
	}

	// Respond with a success status
	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})

}

func handleUserRegistration(userId string) {
    log.Info().Msgf("User with ID %s registered", userId)
    userBalance1, err := db_payment.GetUserBalance(userId)
    if err != nil {
        log.Error().Err(err).Msg("Failed to get user balance")
        return
    }
    log.Info().Msgf("User balance: %v", userBalance1)
    
    db_payment.InsertUserBalance(userId)
    
    userBalance2, err := db_payment.GetUserBalance(userId)
    if err != nil {
        log.Error().Err(err).Msg("Failed to get user balance")
        return
    }
    log.Info().Msgf("User balance: %v", userBalance2)
}

func handleUserDeletion(userId string) {
    log.Info().Msgf("User with ID %s deleted", userId)
    
    userBalance3, err := db_payment.GetUserBalance(userId)
    if err != nil {
        log.Error().Err(err).Msg("Failed to get user balance")
        return
    }
    log.Info().Msgf("User balance: %v", userBalance3)
    
    db_payment.DeleteUserBalance(userId)
    
    // Fixing userId4 typo and userBalance4 unused variable
    userBalance4, err := db_payment.GetUserBalance(userId)
    if err != nil {
        log.Error().Err(err).Msg("Failed to get user balance")
        return
    }
    log.Info().Msgf("User balance: %v", userBalance4)
}

