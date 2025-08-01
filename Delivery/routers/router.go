package routers

import (
	"blog-api/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(blogController *controllers.BlogController) *gin.Engine {
	router := gin.Default()

	router.POST("/blogs", blogController.CreateBlog)
	return router
}
