package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ThrowInternalServerErrorException(c *gin.Context, message string) {
	throwError(c, http.StatusInternalServerError, message)
}

func ThrowUnauthorizedException(c *gin.Context, message string) {
	throwError(c, http.StatusUnauthorized, message)
}

func throwError(c *gin.Context, statuscode int, message string) {
	c.JSON(statuscode, gin.H{"error": message})
	c.Abort()
}
