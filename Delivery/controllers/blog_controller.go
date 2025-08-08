package controllers

import (
	"blog-api/Domain/models"
	"blog-api/usecases"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	blogUC usecases.BlogUseCase
}

func NewBlogController(blogUC usecases.BlogUseCase) *BlogController {
	return &BlogController{blogUC: blogUC}
}

func (ctrl *BlogController) CreateBlog(c *gin.Context) {
	var req CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authorName, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	blog := models.Blog{
		Title:       req.Title,
		Content:     req.Content,
		AuthorID:    userID.(string),
		AuthorName:  authorName.(string),
		Tags:        req.Tags,
		IsPublished: req.IsPublished,
		ViewCount:   0,
		Likes:       0,
		Dislikes:    0,
		Comments:    []models.Comment{},
	}

	createdBlog, err := ctrl.blogUC.CreateBlog(blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ctrl.blogToResponse(createdBlog)
	c.JSON(http.StatusCreated, response)
}

// GET /blogs - Get paginated blogs
func (ctrl *BlogController) GetPaginatedBlogs(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	blogs, err := ctrl.blogUC.GetPaginatedBlogs(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var blogResponses []BlogResponse
	for _, blog := range blogs {
		blogResponses = append(blogResponses, ctrl.blogToResponse(blog))
	}

	response := PaginatedBlogsResponse{
		Blogs:      blogResponses,
		Page:       page,
		Limit:      limit,
		Total:      int64(len(blogs)),
		TotalPages: (len(blogs) + limit - 1) / limit,
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *BlogController) GetBlogByID(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	blog, err := ctrl.blogUC.GetBlogByID(blogID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// Increment view count
	go func() {
		ctrl.blogUC.IncrementViewCount(blogID)
	}()

	response := ctrl.blogToResponse(blog)
	c.JSON(http.StatusOK, response)
}

// PUT /blogs/:id - Update blog
func (ctrl *BlogController) UpdateBlog(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	var req UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing blog
	existingBlog, err := ctrl.blogUC.GetBlogByID(blogID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// Check if user is the author
	userID, exists := c.Get("userID")
	if !exists || userID.(string) != existingBlog.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own blogs"})
		return
	}

	// Update fields
	if req.Title != "" {
		existingBlog.Title = req.Title
	}
	if req.Content != "" {
		existingBlog.Content = req.Content
	}
	if req.Tags != nil {
		existingBlog.Tags = req.Tags
	}
	if req.IsPublished != nil {
		existingBlog.IsPublished = *req.IsPublished
	}

	updatedBlog, err := ctrl.blogUC.UpdateBlog(existingBlog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ctrl.blogToResponse(updatedBlog)
	c.JSON(http.StatusOK, response)
}

// DELETE /blogs/:id - Delete blog
func (ctrl *BlogController) DeleteBlog(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	// Get existing blog to check ownership
	existingBlog, err := ctrl.blogUC.GetBlogByID(blogID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// Check if user is the author
	userID, exists := c.Get("userID")
	if !exists || userID.(string) != existingBlog.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own blogs"})
		return
	}

	err = ctrl.blogUC.DeleteBlog(blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}

// GET /blogs/search - Search blogs
func (ctrl *BlogController) SearchBlogs(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	blogs, err := ctrl.blogUC.SearchBlogs(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var blogResponses []BlogResponse
	for _, blog := range blogs {
		blogResponses = append(blogResponses, ctrl.blogToResponse(blog))
	}

	c.JSON(http.StatusOK, gin.H{"blogs": blogResponses, "query": query})
}

// GET /blogs/filter - Filter blogs
func (ctrl *BlogController) FilterBlogs(c *gin.Context) {
	tags := c.QueryArray("tags")
	dateRange := c.QueryArray("dateRange")
	sortBy := c.DefaultQuery("sortBy", "created_at")

	var dateRangeArray [2]string
	if len(dateRange) >= 2 {
		dateRangeArray[0] = dateRange[0]
		dateRangeArray[1] = dateRange[1]
	}

	blogs, err := ctrl.blogUC.FilterBlogs(tags, dateRangeArray, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var blogResponses []BlogResponse
	for _, blog := range blogs {
		blogResponses = append(blogResponses, ctrl.blogToResponse(blog))
	}

	c.JSON(http.StatusOK, gin.H{"blogs": blogResponses})
}

// POST /blogs/:id/view - Increment view count
func (ctrl *BlogController) IncrementViewCount(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	err := ctrl.blogUC.IncrementViewCount(blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View count incremented"})
}

// POST /blogs/:id/like - Like blog
func (ctrl *BlogController) LikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	err := ctrl.blogUC.UpdateLikes(blogID, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog liked successfully"})
}

// DELETE /blogs/:id/like - Unlike blog
func (ctrl *BlogController) UnlikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	err := ctrl.blogUC.UpdateLikes(blogID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog unliked successfully"})
}

// POST /blogs/:id/dislike - Dislike blog
func (ctrl *BlogController) DislikeBlog(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	err := ctrl.blogUC.UpdateDislikes(blogID, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog disliked successfully"})
}

// DELETE /blogs/:id/dislike - Remove dislike
func (ctrl *BlogController) RemoveDislike(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	err := ctrl.blogUC.UpdateDislikes(blogID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog dislike removed successfully"})
}

// POST /blogs/:id/comments - Add comment
func (ctrl *BlogController) AddComment(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authorName, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	comment := models.Comment{
		BlogID:     blogID,
		AuthorID:   userID.(string),
		AuthorName: authorName.(string),
		Content:    req.Content,
	}

	createdComment, err := ctrl.blogUC.AddComment(blogID, comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ctrl.commentToResponse(createdComment)
	c.JSON(http.StatusCreated, response)
}

// GET /blogs/:id/comments - Get comments
func (ctrl *BlogController) GetComments(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	comments, err := ctrl.blogUC.GetComments(blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var commentResponses []CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, ctrl.commentToResponse(comment))
	}

	c.JSON(http.StatusOK, gin.H{"comments": commentResponses})
}

// Helper methods
func (ctrl *BlogController) blogToResponse(blog models.Blog) BlogResponse {
	var commentResponses []CommentResponse
	for _, comment := range blog.Comments {
		commentResponses = append(commentResponses, ctrl.commentToResponse(comment))
	}

	return BlogResponse{
		ID:          blog.ID,
		Title:       blog.Title,
		Content:     blog.Content,
		AuthorID:    blog.AuthorID,
		AuthorName:  blog.AuthorName,
		Tags:        blog.Tags,
		ViewCount:   blog.ViewCount,
		Likes:       blog.Likes,
		Dislikes:    blog.Dislikes,
		Comments:    commentResponses,
		IsPublished: blog.IsPublished,
		CreatedAt:   blog.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   blog.UpdatedAt.Format(time.RFC3339),
	}
}

func (ctrl *BlogController) commentToResponse(comment models.Comment) CommentResponse {
	return CommentResponse{
		ID:         comment.ID,
		BlogID:     comment.BlogID,
		AuthorID:   comment.AuthorID,
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  comment.UpdatedAt.Format(time.RFC3339),
	}
}
