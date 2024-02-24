package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ThrowInternalServerErrorException(c *gin.Context, message string) {
	throwError(c, http.StatusInternalServerError, message)
}

func ThrowNotFoundException(c *gin.Context, message string) {
	throwError(c, http.StatusNotFound, message)
}

func ThrowUnauthorizedException(c *gin.Context, message string) {
	throwError(c, http.StatusUnauthorized, message)
}

func ThrowBadRequestException(c *gin.Context, message string) {
	throwError(c, http.StatusBadRequest, message)
}

func ThrowServiceUnavailableException(c *gin.Context, message string){
	throwError(c, http.StatusServiceUnavailable, message)
}

func throwError(c *gin.Context, statuscode int, message string) {
	c.JSON(statuscode, ErrorResponse{Error: message})
	c.Abort()
}
