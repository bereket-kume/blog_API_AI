package routers

import (
	"blog-api/Delivery/controllers"
	"blog-api/Delivery/middlewares"
	"blog-api/Domain/interfaces"
	"blog-api/usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, userUC usecases.UserUsecaseInterface, blogUC usecases.BlogUseCase, aiSuggestionUC interfaces.AISuggestionUseCase, tokenService interfaces.TokenService) {
	// Initialize controllers
	userController := controllers.NewUserController(userUC)
	blogController := controllers.NewBlogController(blogUC)
	aiSuggestionController := controllers.NewAISuggestionController(aiSuggestionUC)

	// Initialize profile controllers
	controllers.InitUserController(userUC)

	// Public routes
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/refresh", userController.RefreshToken)

	router.GET("/blogs", blogController.GetPaginatedBlogs)
	router.GET("/blogs/search", blogController.SearchBlogs)
	router.GET("/blogs/filter", blogController.FilterBlogs)
	router.GET("/blogs/:id", blogController.GetBlogByID)
	router.GET("/blogs/:id/comments", blogController.GetComments)

	// Protected routes with real authentication
	auth := router.Group("/api")
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
	}
}
