package routers

import (
	"blog-api/Delivery/controllers"
	"blog-api/Delivery/middlewares"
	"blog-api/Domain/interfaces"
	"blog-api/usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userUC usecases.UserUsecaseInterface, blogUC usecases.BlogUseCase, recommendationUC interfaces.RecommendationUseCase, aiSuggestionUC interfaces.AISuggestionUseCase, tokenService interfaces.TokenService) {
	// Initialize controllers
	userController := controllers.NewUserController(userUC)
	blogController := controllers.NewBlogController(blogUC)
	recommendationController := controllers.NewRecommendationController(recommendationUC)
	aiSuggestionController := controllers.NewAISuggestionController(aiSuggestionUC)

	// Initialize profile controllers
	controllers.InitUserController(userUC)

	// Public routes
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	r.POST("/refresh", userController.RefreshToken)
	r.GET("/verify-email", userController.VerifyEmail)
	r.POST("/forgot-password", userController.RequestPasswordReset)
	r.GET("/reset-password", userController.ResetPassword)

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
		// User profile routes with real auth
		user := auth.Group("/user").Use(middlewares.AuthMiddleware(tokenService))
		{
			user.GET("/profile", controllers.GetUserProfile)
			user.PUT("/profile", controllers.UpdateUserProfile)
		}

		// Blog routes with real auth
		blogs := auth.Group("/blogs").Use(middlewares.AuthMiddleware(tokenService))
		{
			blogs.POST("/", blogController.CreateBlog)
			blogs.PUT("/:id", blogController.UpdateBlog)
			blogs.DELETE("/:id", blogController.DeleteBlog)
			blogs.POST("/:id/comments", blogController.AddComment)
			blogs.POST("/:id/like", blogController.LikeBlog)
			blogs.POST("/:id/unlike", blogController.UnlikeBlog)
			blogs.POST("/:id/dislike", blogController.DislikeBlog)
			blogs.POST("/:id/remove-dislike", blogController.RemoveDislike)
		}

		// AI routes with real auth
		ai := auth.Group("/ai").Use(middlewares.AuthMiddleware(tokenService))
		{
			ai.POST("/suggestions", controllers.GenerateAISuggestion)
			ai.POST("/ideas", controllers.GenerateContentIdeas)
			ai.POST("/save", aiSuggestionController.SaveAISuggestion)
			ai.GET("/suggestions", aiSuggestionController.GetAISuggestions)
			ai.GET("/suggestions/status/:status", aiSuggestionController.GetAISuggestionsByStatus)
			ai.POST("/suggestions/:id/convert-to-draft", aiSuggestionController.ConvertSuggestionToDraft)
			ai.DELETE("/suggestions/:id", aiSuggestionController.DeleteAISuggestion)
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
		admin := auth.Group("/admin").Use(middlewares.AuthMiddleware(tokenService, "admin", "superadmin"))
		{
			admin.POST("/promote", userController.Promote)
		}

		// Superadmin-only routes
		superadmin := auth.Group("/superadmin").Use(middlewares.AuthMiddleware(tokenService, "superadmin"))
		{
			superadmin.POST("/demote", userController.Demote)
		}

		// Logout route (all authenticated users)
		auth.POST("/logout", middlewares.AuthMiddleware(tokenService), userController.Logout)
	}
}
