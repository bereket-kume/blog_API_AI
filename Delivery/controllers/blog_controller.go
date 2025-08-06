package controllers

import (
	domain "blog-api/Domain/models"
	usecases "blog-api/Usecases"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Extract user ID from context (assuming it's set by auth middleware)
	userIDStr := ctx.GetString("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	blog.AuthorID = userID

	createdBlog, err := c.blogUseCase.CreateBlog(blog)
	if err != nil {
		log.Printf("error creating blog: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdBlog)
}

func (c *BlogController) GetAllBlogs(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	blogs, err := c.blogUseCase.GetPaginatedBlogs(page, limit)
	if err != nil {
		log.Printf("error getting blogs: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get blogs", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, blogs)
}

func (c *BlogController) GetBlogByID(ctx *gin.Context) {
	blogID := ctx.Param("id")
	if blogID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	// Get the blog
	blog, err := c.blogUseCase.GetBlogByID(blogID)
	if err != nil {
		log.Printf("error getting blog by ID: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog not found", "details": err.Error()})
		return
	}

	// Increment view count
	if err := c.blogUseCase.IncrementViewCount(blogID); err != nil {
		log.Printf("error incrementing view count: %v", err)
	}

	ctx.JSON(http.StatusOK, blog)
}

func (c *BlogController) UpdateBlog(ctx *gin.Context) {
	blogID := ctx.Param("id")
	if blogID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	var blog domain.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Parse blog ID
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}
	blog.ID = objectID

	updatedBlog, err := c.blogUseCase.UpdateBlog(blog)
	if err != nil {
		log.Printf("error updating blog: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedBlog)
}

func (c *BlogController) DeleteBlog(ctx *gin.Context) {
	blogID := ctx.Param("id")
	if blogID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	err := c.blogUseCase.DeleteBlog(blogID)
	if err != nil {
		log.Printf("error deleting blog: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}

func (c *BlogController) SearchBlogs(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	blogs, err := c.blogUseCase.SearchBlogs(query)
	if err != nil {
		log.Printf("error searching blogs: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search blogs", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, blogs)
}

func (c *BlogController) FilterBlogs(ctx *gin.Context) {
	tags := ctx.QueryArray("tags")
	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")

	var dateRange [2]string
	if dateFrom != "" && dateTo != "" {
		dateRange = [2]string{dateFrom, dateTo}
	}

	blogs, err := c.blogUseCase.FilterBlogs(tags, dateRange, sortBy)
	if err != nil {
		log.Printf("error filtering blogs: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter blogs", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, blogs)
}

func (c *BlogController) UpdateLikes(ctx *gin.Context) {
	blogID := ctx.Param("id")
	if blogID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	var req struct {
		Increment bool `json:"increment"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	err := c.blogUseCase.UpdateLikes(blogID, req.Increment)
	if err != nil {
		log.Printf("error updating likes: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update likes", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Likes updated successfully"})
}

func (c *BlogController) UpdateDislikes(ctx *gin.Context) {
	blogID := ctx.Param("id")
	if blogID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	var req struct {
		Increment bool `json:"increment"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	err := c.blogUseCase.UpdateDislikes(blogID, req.Increment)
	if err != nil {
		log.Printf("error updating dislikes: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update dislikes", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Dislikes updated successfully"})
}

func (c *BlogController) AddComment(ctx *gin.Context) {
	blogID := ctx.Param("id")
	if blogID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	var comment domain.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Extract user ID from context
	userIDStr := ctx.GetString("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	comment.UserID = userID

	createdComment, err := c.blogUseCase.AddComment(blogID, comment)
	if err != nil {
		log.Printf("error adding comment: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdComment)
}

func (c *BlogController) GetComments(ctx *gin.Context) {
	blogID := ctx.Param("id")
	if blogID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	comments, err := c.blogUseCase.GetComments(blogID)
	if err != nil {
		log.Printf("error getting comments: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
