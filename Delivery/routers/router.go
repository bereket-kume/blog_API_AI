package routers

import (
	"blog-api/Delivery/controllers"
	Database "blog-api/Infrastructure/database"
	"blog-api/Infrastructure/middleware"
	"blog-api/Infrastructure/repositories"
	usecases "blog-api/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Initialize User Dependencies
	userCollection := Database.GetUserCollection()
	userRepo := repositories.NewUserRepository(userCollection)

	userUsecase := usecases.NewUserUsecase(userRepo)
	controllers.InitUserController(userUsecase)

	// User Routes (Protected)
	user := router.Group("/user")
	user.Use(middleware.MockAuth()) // Replace with real auth middleware when ready
	{
		user.GET("/profile", controllers.GetUserProfile)
		user.PUT("/profile", controllers.UpdateUserProfile)
	}

	// AI Suggestion Routes (Also protected — mock for now)
	ai := router.Group("/ai")
	ai.Use(middleware.MockAuth())
	{
		ai.POST("/generate", controllers.GenerateAISuggestion)
	}

	return router
}
