package utils

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SendSuccess sends a standardized success response
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// SendError sends a standardized error response
func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}
