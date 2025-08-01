package controllers

import (
	domain "blog-api/Domain/models"
	usecases "blog-api/Usecases"
	"log"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	blogUseCase usecases.BlogUseCase
}

func NewBlogController(blogUseCase usecases.BlogUseCase) *BlogController {
	return &BlogController{
		blogUseCase: blogUseCase,
	}
}

func (c *BlogController) CreateBlog(ctx *gin.Context) {
	var blog domain.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	// userID, err := primitive.ObjectIDFromHex(ctx.GetHeader("user_id"))
	// if err != nil {
	// 	ctx.JSON(400, gin.H{"error": "Invalid user ID"})
	// 	return
	// }
	// blog.AuthorID = userID
	createdBlog, err := c.blogUseCase.CreateBlog(blog)
	log.Printf("error creating blog: %v", err)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create blog", "details": err.Error()})
		return
	}
	ctx.JSON(201, createdBlog)
}

func (c *BlogController) GetAllBlogs(ctx *gin.Context) {

}
