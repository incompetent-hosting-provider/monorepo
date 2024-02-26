package payment

import (
	"incompetent-hosting-provider/backend/pkg/constants"
	db_payment "incompetent-hosting-provider/backend/pkg/db/tables/payment"
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type BalanceResponse struct {
	Balance int `json:"balance"`
}

type BalanceIncreaseBody struct {
	Amount int `json:"amount"`
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
		util.ThrowInternalServerErrorException(c, "Internal Server Error")
		return
	}
	c.JSON(http.StatusOK, BalanceResponse{Balance: balance})
}

// godoc
// @Summary 					Increase user balance
//
// @Schemes
// @Description 				Increase user balance
// @Tags 						payment
// @Produce 					json
//
// @Security					BearerAuth
//
// @Param request body payment.BalanceIncreaseBody true "query params"
//
// @Success 					200 {object} payment.BalanceResponse
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						404 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /payment [post]
func ChangeCreditHandler(c *gin.Context) {

	var updateRequest BalanceIncreaseBody

	err := c.ShouldBindJSON(&updateRequest)
	if err != nil {
		log.Info().Msg("Could not parse request body")
		util.ThrowBadRequestException(c, "Could not parse body.")
		return
	}

	// Use header set by middleware
	userId := c.GetHeader(constants.USER_ID_HEADER)

	newBalance, err := db_payment.IncreaseBalance(userId, updateRequest.Amount)
	if err != nil {
		log.Warn().Msg("Error in DB write")
		util.ThrowInternalServerErrorException(c, "Internal Server Error")
		return
	}
	c.JSON(http.StatusCreated, BalanceResponse{Balance: newBalance})
}
