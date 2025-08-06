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

	// Protected routes
	auth := r.Group("/api")
	{

		// Admin-only routes
		admin := auth.Group("/admin").Use(middlewares.AuthMiddleware(tokenService, "admin"))
		{
			admin.POST("/promote/:email", userController.Promote)
		}

		// Superadmin-only routes
		super := auth.Group("/superadmin").Use(middlewares.AuthMiddleware(tokenService, "superadmin"))
		{
			super.POST("/demote/:email", userController.Demote)
		}
	}
}
