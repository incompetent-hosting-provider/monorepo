package payment

import (
	"incompetent-hosting-provider/backend/pkg/constants"
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type BalanceResponse struct {
	Balance int `json:"balance"`
}

// godoc
// @Summary 					Get user balance
//
// @Schemes
// @Description 				Get current user balance
// @Tags 						payment
// @Produce 					json
//
// @Security					BearerAuth
//
// @Success 					200 {object} payment.BalanceResponse
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						404 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /payment [get]
func CreditFetchHandler(c *gin.Context) {
	// Use header set by middleware
	userId := c.GetHeader(constants.USER_ID_HEADER)
	balance, err := GetCurrentCredits(userId)
	if err != nil {
		log.Warn().Msgf("Tried to fetch account balance but failed due to an error %s", err)
		util.ThrowInternalServerErrorException(c, "Could not fetch balance. This is not an authorization issue.")
		return
	}
	c.JSON(http.StatusOK, BalanceResponse{Balance: balance})
}

func GetCurrentCredits(userId string) (int, error) {
	return 1000, nil
}
