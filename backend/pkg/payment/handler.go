package payment

import (
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CreditFetchHandler(c *gin.Context) {
	// Use header set by middleware
	userId := c.GetHeader("user-id")
	balance, err := getCurrentCredits(userId)
	if err != nil {
		log.Warn().Msgf("Tried to fetch account balance but failed due to an error %s", err)
		util.ThrowInternalServerErrorException(c, "Could not fetch balance. This is not an authorization issue.")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

func getCurrentCredits(userId string) (int, error) {
	return 1000, nil
}
