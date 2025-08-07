package routers

import (
	"blog-api/Delivery/controllers"
	"blog-api/Infrastructure/middleware"
	"blog-api/Infrastructure/repositories"
	usecases "blog-api/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	userRepo := repositories.NewUserRepository()
	userUsecase := usecases.NewUserUsecase(userRepo)
	controllers.InitUserController(userUsecase)

	user := router.Group("/user")
	user.Use(middleware.MockAuth())
	{
		user.GET("/profile", controllers.GetUserProfile)
		user.PUT("/profile", controllers.UpdateUserProfile)
	}

	return router
}
