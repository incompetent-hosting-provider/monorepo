package webhook

import (
	db_instances "incompetent-hosting-provider/backend/pkg/db/tables/instances"
	db_payment "incompetent-hosting-provider/backend/pkg/db/tables/payment"
	"incompetent-hosting-provider/backend/pkg/mq_handler"
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type KeycloakWebookPayload struct {
	// This is either DELETE_ACCOUNT and REGISTER
	KeycloakEvent string `json:"keycloakEvent"`
	UserSub       string `json:"userId"`
}

// godoc
// @Summary 					SPI webhook handler
//
// @Schemes
// @Description 				Handle keycloak SPI webhook events
// @Tags 						webhook
//
// @Success 					200 {string} string	"ok"
//
// @Failure						500 {object} util.ErrorResponse
//
// @Router /spi-webhook [post]
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
		handleUserRegistration(receivedEvent.UserSub)
	case "DELETE_ACCOUNT":
		handleUserDeletion(receivedEvent.UserSub)
	default:
		log.Warn().Msgf("Unsupported Keycloak event type: %s", receivedEvent.KeycloakEvent)
	}

	// Respond with a success status
	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func handleUserRegistration(userSub string) {
	err := db_payment.InsertUserBalance(userSub) // Insert user balance

	if err != nil {
		log.Warn().Msgf("Could not create user balance due to an error: %v", err)
	}
}

func handleUserDeletion(userSub string) {
	err := db_payment.DeleteUserBalance(userSub) // Delete user balance

	if err != nil {
		log.Warn().Msgf("Could not delete user balance: %v", err)
	}

	userInstances, err := db_instances.GetAllUserInstances(userSub)

	if err != nil {
		log.Warn().Msgf("Could not fetch user instances due to an error: %v", err)
	}

	for _, v := range userInstances {
		err := mq_handler.PublishDeleteContainerEvent(mq_handler.DeleteContainerEvent{
			UserId:        userSub,
			ContainerUUID: v.ContainerUUID,
		})

		if err != nil {
			log.Warn().Msgf("Could not delete user instance due to an error: %v", err)
		}

		err = db_instances.DeleteInstanceById(userSub, v.ContainerUUID)

		if err != nil {
			log.Warn().Msgf("Could not delete user instance due to an error: %v", err)
		}
		log.Debug().Msg("Deleted and stopped user instance")
	}
}
