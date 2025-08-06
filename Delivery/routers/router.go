package routers

import (
	"blog-api/Delivery/controllers"
	"blog-api/Delivery/middlewares"
	"blog-api/Domain/interfaces"
	"blog-api/usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userUC usecases.UserUsecaseInterface, tokenService interfaces.TokenService) {
	// Controllers
	userController := controllers.NewUserController(userUC)

	// Public routes
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	r.POST("/refresh", userController.RefreshToken)

	// Protected routes under /api
	api := r.Group("/api")
	{
		// Promote (admin or superadmin)
		api.POST("/promote/:email",
			middlewares.AuthMiddleware(tokenService, "admin", "superadmin"),
			userController.Promote)

		// Demote (superadmin only)
		api.POST("/demote/:email",
			middlewares.AuthMiddleware(tokenService, "superadmin"),
			userController.Demote)
	}
}
