// Inside your middleware/mock.go file
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MockAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use an email that DOES NOT exist in your database for this test.
		c.Set("user_id", primitive.NewObjectID())
		c.Set("email", "nonexistentuser@example.com")
		c.Next()
	}
}
