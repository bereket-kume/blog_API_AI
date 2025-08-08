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
