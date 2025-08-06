package middleware

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MockAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := primitive.ObjectIDFromHex("6893b95f2f0bd5cf28b04d01")
		c.Set("user_id", userID)
		c.Set("email", "realuser@example.com") // Replace with the actual email in your DB
		c.Next()
	}
}
