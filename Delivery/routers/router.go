package routers

import (
	"blog-api/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(blogController *controllers.BlogController) *gin.Engine {
	router := gin.Default()

	// Blog routes
	blogs := router.Group("/blogs")
	{
		blogs.POST("/", blogController.CreateBlog)                 // Create a new blog
		blogs.GET("/", blogController.GetAllBlogs)                 // Get all blogs with pagination
		blogs.GET("/search", blogController.SearchBlogs)           // Search blogs
		blogs.GET("/filter", blogController.FilterBlogs)           // Filter blogs by tags, date, etc.
		blogs.GET("/:id", blogController.GetBlogByID)              // Get blog by ID
		blogs.PUT("/:id", blogController.UpdateBlog)               // Update blog
		blogs.DELETE("/:id", blogController.DeleteBlog)            // Delete blog
		blogs.POST("/:id/likes", blogController.UpdateLikes)       // Update likes
		blogs.POST("/:id/dislikes", blogController.UpdateDislikes) // Update dislikes
		blogs.POST("/:id/comments", blogController.AddComment)     // Add comment to blog
		blogs.GET("/:id/comments", blogController.GetComments)     // Get comments for blog
	}

	return router
}
