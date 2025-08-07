package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standardized success response
func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// Standardized error response
func SendError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"message": message,
	})
}
