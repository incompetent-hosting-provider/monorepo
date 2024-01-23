package user

import (
	"incompetent-hosting-provider/backend/pkg/constants"
	"incompetent-hosting-provider/backend/pkg/payment"
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	Balance int    `json:"balance"`
	Email   string `json:"email"`
}

// godoc
// @Summary 					Get user info
//
// @Schemes
// @Description 				Get user info
// @Tags 						user
// @Produce 					json
//
// @Security					BearerAuth
//
// @Success 					200 {object} user.UserResponse
//
// @Failure						401 {object} util.ErrorResponse
// @Failure						404 {object} util.ErrorResponse
// @Failure						500 {object} util.ErrorResponse
//
// @Router /user [get]
func UserFetchHandler(c *gin.Context) {
	// Use header set by middleware
	userMail := c.Request.Header.Get(constants.USER_EMAIL_HEADER)
	userId := c.Request.Header.Get(constants.USER_ID_HEADER)

	balance, err := payment.GetCurrentCredits(userId)

	if err != nil {
		util.ThrowInternalServerErrorException(c, "Could not fetch balance")
	}

	resp := UserResponse{
		Email:   userMail,
		Balance: balance,
	}

	c.JSON(http.StatusOK, resp)
}
