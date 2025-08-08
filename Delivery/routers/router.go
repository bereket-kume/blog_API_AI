package routers

import (
	"blog-api/Delivery/controllers"
	"blog-api/Delivery/middlewares"
	"blog-api/Domain/interfaces"
	"blog-api/usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userUC usecases.UserUsecaseInterface, blogUC usecases.BlogUseCase, recommendationUC interfaces.RecommendationUseCase, tokenService interfaces.TokenService) {
	// Controllers
	userController := controllers.NewUserController(userUC)
	blogController := controllers.NewBlogController(blogUC)
	recommendationController := controllers.NewRecommendationController(recommendationUC)

	// Initialize profile controller
	controllers.InitUserController(userUC)

	// Public routes
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	r.POST("/refresh", userController.RefreshToken)
	r.GET("/verify-email", userController.VerifyEmail)
	r.POST("/forgot-password", userController.RequestPasswordReset)
	r.GET("/reset-password", userController.ResetPassword)

	// Blog routes (public)
	r.GET("/blogs", blogController.GetPaginatedBlogs)
	r.GET("/blogs/search", blogController.SearchBlogs)
	r.GET("/blogs/filter", blogController.FilterBlogs)
	r.GET("/blogs/:id", blogController.GetBlogByID)
	r.GET("/blogs/:id/comments", blogController.GetComments)

	// Recommendation routes (public)
	r.GET("/recommendations/trending", recommendationController.GetTrendingContent)
	r.GET("/recommendations/popular", recommendationController.GetPopularContent)
	r.GET("/recommendations/new", recommendationController.GetNewContent)
	r.GET("/recommendations/discovery", recommendationController.GetContentDiscovery)
	r.GET("/blogs/:id/similar", recommendationController.GetSimilarContent)

	// Protected routes
	auth := r.Group("/api")
	{
		// Blog routes (authenticated)
		blogs := auth.Group("/blogs").Use(middlewares.AuthMiddleware(tokenService))
		{
			blogs.POST("/", blogController.CreateBlog)
			blogs.PUT("/:id", blogController.UpdateBlog)
			blogs.DELETE("/:id", blogController.DeleteBlog)
			blogs.POST("/:id/comments", blogController.AddComment)
			blogs.POST("/:id/view", blogController.IncrementViewCount)
			blogs.POST("/:id/like", blogController.LikeBlog)
			blogs.DELETE("/:id/like", blogController.UnlikeBlog)
			blogs.POST("/:id/dislike", blogController.DislikeBlog)
			blogs.DELETE("/:id/dislike", blogController.RemoveDislike)
		}

		// Profile routes (authenticated)
		profile := auth.Group("/profile").Use(middlewares.AuthMiddleware(tokenService))
		{
			profile.GET("/", controllers.GetUserProfile)
			profile.PUT("/", controllers.UpdateUserProfile)
		}

		// Recommendation routes (authenticated)
		recommendations := auth.Group("/recommendations").Use(middlewares.AuthMiddleware(tokenService))
		{
			recommendations.POST("/track", recommendationController.TrackUserAction)
			recommendations.GET("/personal", recommendationController.GetUserRecommendations)
			recommendations.GET("/interests", recommendationController.GetUserInterests)
			recommendations.GET("/behavior", recommendationController.GetUserBehaviorSummary)
			recommendations.GET("/stats", recommendationController.GetRecommendationStats)
			recommendations.PUT("/:id/view", recommendationController.MarkRecommendationViewed)
		}

		// Admin-only routes
		auth.POST("/promote/",
			middlewares.AuthMiddleware(tokenService, "admin", "superadmin"),
			userController.Promote)

		// Demote (superadmin only)
		auth.POST("/demote/",
			middlewares.AuthMiddleware(tokenService, "superadmin"),
			userController.Demote)
		auth.POST("/logout", middlewares.AuthMiddleware(tokenService, "user", "admin", "superadmin"), userController.Logout)
	}
}
