package controllers

import (
	"blog-api/Domain/models"
	"blog-api/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BlogControllerTestSuite struct {
	suite.Suite
	router     *gin.Engine
	controller *BlogController
	mockUC     *mocks.BlogUseCaseMock
}

func (suite *BlogControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.mockUC = &mocks.BlogUseCaseMock{}
	suite.controller = NewBlogController(suite.mockUC)
}

func (suite *BlogControllerTestSuite) TearDownTest() {
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *BlogControllerTestSuite) TestCreateBlog_Success() {
	// Test data
	requestBody := CreateBlogRequest{
		Title:       "Test Blog",
		Content:     "Test Content",
		Tags:        []string{"test", "go"},
		IsPublished: true,
	}
	jsonBody, _ := json.Marshal(requestBody)

	expectedBlog := models.Blog{
		ID:          "blog123",
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Tags:        []string{"test", "go"},
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Setup mock
	suite.mockUC.On("CreateBlog", mock.AnythingOfType("models.Blog")).Return(expectedBlog, nil)

	// Setup route
	suite.router.POST("/blogs", func(c *gin.Context) {
		c.Set("userID", "user123")
		c.Set("email", "test@example.com")
		suite.controller.CreateBlog(c)
	})

	// Create request
	req, _ := http.NewRequest("POST", "/blogs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response BlogResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedBlog.ID, response.ID)
	assert.Equal(suite.T(), expectedBlog.Title, response.Title)
	assert.Equal(suite.T(), expectedBlog.Content, response.Content)
}

func (suite *BlogControllerTestSuite) TestCreateBlog_InvalidRequest() {
	// Test data with missing required fields
	requestBody := map[string]interface{}{
		"content": "Test Content",
		// Missing title
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Setup route
	suite.router.POST("/blogs", func(c *gin.Context) {
		c.Set("userID", "user123")
		c.Set("email", "test@example.com")
		suite.controller.CreateBlog(c)
	})

	// Create request
	req, _ := http.NewRequest("POST", "/blogs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *BlogControllerTestSuite) TestCreateBlog_Unauthorized() {
	// Test data
	requestBody := CreateBlogRequest{
		Title:       "Test Blog",
		Content:     "Test Content",
		Tags:        []string{"test", "go"},
		IsPublished: true,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Setup route without user context
	suite.router.POST("/blogs", suite.controller.CreateBlog)

	// Create request
	req, _ := http.NewRequest("POST", "/blogs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *BlogControllerTestSuite) TestCreateBlog_UseCaseError() {
	// Test data
	requestBody := CreateBlogRequest{
		Title:       "Test Blog",
		Content:     "Test Content",
		Tags:        []string{"test", "go"},
		IsPublished: true,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Setup mock to return error
	suite.mockUC.On("CreateBlog", mock.AnythingOfType("models.Blog")).Return(models.Blog{}, errors.New("database error"))

	// Setup route
	suite.router.POST("/blogs", func(c *gin.Context) {
		c.Set("userID", "user123")
		c.Set("email", "test@example.com")
		suite.controller.CreateBlog(c)
	})

	// Create request
	req, _ := http.NewRequest("POST", "/blogs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *BlogControllerTestSuite) TestGetPaginatedBlogs_Success() {
	// Test data
	expectedBlogs := []models.Blog{
		{
			ID:          "blog1",
			Title:       "Blog 1",
			Content:     "Content 1",
			AuthorID:    "user1",
			AuthorName:  "user1@test.com",
			IsPublished: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "blog2",
			Title:       "Blog 2",
			Content:     "Content 2",
			AuthorID:    "user2",
			AuthorName:  "user2@test.com",
			IsPublished: true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Setup mock
	suite.mockUC.On("GetPaginatedBlogs", 1, 10).Return(expectedBlogs, nil)

	// Setup route
	suite.router.GET("/blogs", suite.controller.GetPaginatedBlogs)

	// Create request
	req, _ := http.NewRequest("GET", "/blogs?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response PaginatedBlogsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response.Blogs, 2)
	assert.Equal(suite.T(), 1, response.Page)
	assert.Equal(suite.T(), 10, response.Limit)
}

func (suite *BlogControllerTestSuite) TestGetPaginatedBlogs_UseCaseError() {
	// Setup mock to return error
	suite.mockUC.On("GetPaginatedBlogs", 1, 10).Return(nil, errors.New("database error"))

	// Setup route
	suite.router.GET("/blogs", suite.controller.GetPaginatedBlogs)

	// Create request
	req, _ := http.NewRequest("GET", "/blogs?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *BlogControllerTestSuite) TestGetBlogByID_Success() {
	// Test data
	expectedBlog := models.Blog{
		ID:          "blog123",
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Setup mock
	suite.mockUC.On("GetBlogByID", "blog123").Return(expectedBlog, nil)
	// IncrementViewCount is called in a goroutine, so we expect it but don't require it
	suite.mockUC.On("IncrementViewCount", "blog123").Return(nil).Maybe()

	// Setup route
	suite.router.GET("/blogs/:id", suite.controller.GetBlogByID)

	// Create request
	req, _ := http.NewRequest("GET", "/blogs/blog123", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response BlogResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedBlog.ID, response.ID)
	assert.Equal(suite.T(), expectedBlog.Title, response.Title)
}

func (suite *BlogControllerTestSuite) TestGetBlogByID_NotFound() {
	// Setup mock to return error
	suite.mockUC.On("GetBlogByID", "nonexistent").Return(models.Blog{}, errors.New("blog not found"))

	// Setup route
	suite.router.GET("/blogs/:id", suite.controller.GetBlogByID)

	// Create request
	req, _ := http.NewRequest("GET", "/blogs/nonexistent", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *BlogControllerTestSuite) TestUpdateBlog_Success() {
	// Test data
	requestBody := UpdateBlogRequest{
		Title:       "Updated Blog",
		Content:     "Updated Content",
		Tags:        []string{"updated", "go"},
		IsPublished: boolPtr(true),
	}
	jsonBody, _ := json.Marshal(requestBody)

	existingBlog := models.Blog{
		ID:          "blog123",
		Title:       "Original Blog",
		Content:     "Original Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updatedBlog := models.Blog{
		ID:          "blog123",
		Title:       "Updated Blog",
		Content:     "Updated Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Tags:        []string{"updated", "go"},
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Setup mock
	suite.mockUC.On("GetBlogByID", "blog123").Return(existingBlog, nil)
	suite.mockUC.On("UpdateBlog", mock.AnythingOfType("models.Blog")).Return(updatedBlog, nil)

	// Setup route
	suite.router.PUT("/blogs/:id", func(c *gin.Context) {
		c.Set("userID", "user123")
		suite.controller.UpdateBlog(c)
	})

	// Create request
	req, _ := http.NewRequest("PUT", "/blogs/blog123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response BlogResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Blog", response.Title)
	assert.Equal(suite.T(), "Updated Content", response.Content)
}

func (suite *BlogControllerTestSuite) TestUpdateBlog_Forbidden() {
	// Test data
	requestBody := UpdateBlogRequest{
		Title: "Updated Blog",
	}
	jsonBody, _ := json.Marshal(requestBody)

	existingBlog := models.Blog{
		ID:       "blog123",
		Title:    "Original Blog",
		AuthorID: "different-user", // Different author
	}

	// Setup mock
	suite.mockUC.On("GetBlogByID", "blog123").Return(existingBlog, nil)

	// Setup route
	suite.router.PUT("/blogs/:id", func(c *gin.Context) {
		c.Set("userID", "user123") // Different user trying to update
		suite.controller.UpdateBlog(c)
	})

	// Create request
	req, _ := http.NewRequest("PUT", "/blogs/blog123", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *BlogControllerTestSuite) TestDeleteBlog_Success() {
	// Test data
	existingBlog := models.Blog{
		ID:       "blog123",
		Title:    "Test Blog",
		AuthorID: "user123",
	}

	// Setup mock
	suite.mockUC.On("GetBlogByID", "blog123").Return(existingBlog, nil)
	suite.mockUC.On("DeleteBlog", "blog123").Return(nil)

	// Setup route
	suite.router.DELETE("/blogs/:id", func(c *gin.Context) {
		c.Set("userID", "user123")
		suite.controller.DeleteBlog(c)
	})

	// Create request
	req, _ := http.NewRequest("DELETE", "/blogs/blog123", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BlogControllerTestSuite) TestDeleteBlog_Forbidden() {
	// Test data
	existingBlog := models.Blog{
		ID:       "blog123",
		Title:    "Test Blog",
		AuthorID: "different-user", // Different author
	}

	// Setup mock
	suite.mockUC.On("GetBlogByID", "blog123").Return(existingBlog, nil)

	// Setup route
	suite.router.DELETE("/blogs/:id", func(c *gin.Context) {
		c.Set("userID", "user123") // Different user trying to delete
		suite.controller.DeleteBlog(c)
	})

	// Create request
	req, _ := http.NewRequest("DELETE", "/blogs/blog123", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *BlogControllerTestSuite) TestSearchBlogs_Success() {
	// Test data
	expectedBlogs := []models.Blog{
		{
			ID:          "blog1",
			Title:       "Go Programming",
			Content:     "Learn Go programming",
			AuthorID:    "user1",
			AuthorName:  "user1@test.com",
			IsPublished: true,
		},
		{
			ID:          "blog2",
			Title:       "Web Development with Go",
			Content:     "Learn web development",
			AuthorID:    "user2",
			AuthorName:  "user2@test.com",
			IsPublished: true,
		},
	}

	// Setup mock
	suite.mockUC.On("SearchBlogs", "Go").Return(expectedBlogs, nil)

	// Setup route
	suite.router.GET("/blogs/search", suite.controller.SearchBlogs)

	// Create request
	req, _ := http.NewRequest("GET", "/blogs/search?q=Go", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Go", response["query"])
}

func (suite *BlogControllerTestSuite) TestSearchBlogs_MissingQuery() {
	// Setup route
	suite.router.GET("/blogs/search", suite.controller.SearchBlogs)

	// Create request without query parameter
	req, _ := http.NewRequest("GET", "/blogs/search", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *BlogControllerTestSuite) TestFilterBlogs_Success() {
	// Test data
	expectedBlogs := []models.Blog{
		{
			ID:          "blog1",
			Title:       "Go Blog",
			Content:     "Go content",
			AuthorID:    "user1",
			AuthorName:  "user1@test.com",
			Tags:        []string{"go", "programming"},
			IsPublished: true,
		},
	}

	// Setup mock
	suite.mockUC.On("FilterBlogs", []string{"go"}, [2]string{"2023-01-01", "2023-12-31"}, "created_at").Return(expectedBlogs, nil)

	// Setup route
	suite.router.GET("/blogs/filter", suite.controller.FilterBlogs)

	// Create request
	req, _ := http.NewRequest("GET", "/blogs/filter?tags=go&dateRange=2023-01-01&dateRange=2023-12-31&sortBy=created_at", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response["blogs"])
}

func (suite *BlogControllerTestSuite) TestIncrementViewCount_Success() {
	// Setup mock
	suite.mockUC.On("IncrementViewCount", "blog123").Return(nil)

	// Setup route
	suite.router.POST("/blogs/:id/view", suite.controller.IncrementViewCount)

	// Create request
	req, _ := http.NewRequest("POST", "/blogs/blog123/view", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BlogControllerTestSuite) TestLikeBlog_Success() {
	// Setup mock
	suite.mockUC.On("UpdateLikes", "blog123", true).Return(nil)

	// Setup route
	suite.router.POST("/blogs/:id/like", suite.controller.LikeBlog)

	// Create request
	req, _ := http.NewRequest("POST", "/blogs/blog123/like", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BlogControllerTestSuite) TestUnlikeBlog_Success() {
	// Setup mock
	suite.mockUC.On("UpdateLikes", "blog123", false).Return(nil)

	// Setup route
	suite.router.DELETE("/blogs/:id/like", suite.controller.UnlikeBlog)

	// Create request
	req, _ := http.NewRequest("DELETE", "/blogs/blog123/like", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BlogControllerTestSuite) TestDislikeBlog_Success() {
	// Setup mock
	suite.mockUC.On("UpdateDislikes", "blog123", true).Return(nil)

	// Setup route
	suite.router.POST("/blogs/:id/dislike", suite.controller.DislikeBlog)

	// Create request
	req, _ := http.NewRequest("POST", "/blogs/blog123/dislike", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BlogControllerTestSuite) TestRemoveDislike_Success() {
	// Setup mock
	suite.mockUC.On("UpdateDislikes", "blog123", false).Return(nil)

	// Setup route
	suite.router.DELETE("/blogs/:id/dislike", suite.controller.RemoveDislike)

	// Create request
	req, _ := http.NewRequest("DELETE", "/blogs/blog123/dislike", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BlogControllerTestSuite) TestAddComment_Success() {
	// Test data
	requestBody := CommentRequest{
		Content: "Great blog post!",
	}
	jsonBody, _ := json.Marshal(requestBody)

	expectedComment := models.Comment{
		ID:         "comment123",
		BlogID:     "blog123",
		AuthorID:   "user123",
		AuthorName: "test@example.com",
		Content:    "Great blog post!",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Setup mock
	suite.mockUC.On("AddComment", "blog123", mock.AnythingOfType("models.Comment")).Return(expectedComment, nil)

	// Setup route
	suite.router.POST("/blogs/:id/comments", func(c *gin.Context) {
		c.Set("userID", "user123")
		c.Set("email", "test@example.com")
		suite.controller.AddComment(c)
	})

	// Create request
	req, _ := http.NewRequest("POST", "/blogs/blog123/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response CommentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedComment.ID, response.ID)
	assert.Equal(suite.T(), expectedComment.Content, response.Content)
}

func (suite *BlogControllerTestSuite) TestGetComments_Success() {
	// Test data
	expectedComments := []models.Comment{
		{
			ID:         "comment1",
			BlogID:     "blog123",
			AuthorID:   "user1",
			AuthorName: "user1@test.com",
			Content:    "Comment 1",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "comment2",
			BlogID:     "blog123",
			AuthorID:   "user2",
			AuthorName: "user2@test.com",
			Content:    "Comment 2",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	// Setup mock
	suite.mockUC.On("GetComments", "blog123").Return(expectedComments, nil)

	// Setup route
	suite.router.GET("/blogs/:id/comments", suite.controller.GetComments)

	// Create request
	req, _ := http.NewRequest("GET", "/blogs/blog123/comments", nil)
	w := httptest.NewRecorder()

	// Execute request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response["comments"])
}

// Helper function to create bool pointer
func boolPtr(b bool) *bool {
	return &b
}

// Run the test suite
func TestBlogControllerTestSuite(t *testing.T) {
	suite.Run(t, new(BlogControllerTestSuite))
}
