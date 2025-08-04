package routers

import (
	"blog-api/Delivery/controllers"
	"blog-api/Infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	user := router.Group("/user")
	user.Use(middleware.MockAuth()) // replace with real auth later
	{
		user.GET("/profile", controllers.GetUserProfile)
		user.PUT("/profile", controllers.UpdateUserProfile)
	}

	return router
}
