package middleware

import (
	"github.com/gin-gonic/gin"
)

func MockAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mock authentication logic
		c.Set("userID", "mockUserID")
		c.Next()
	}
}
