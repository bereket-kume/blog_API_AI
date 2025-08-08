package main

import (
	"blog-api/Domain/models"
	"fmt"
	"time"
)

// TestHelper provides common test utilities
type TestHelper struct{}

// CreateTestBlog creates a sample blog for testing
func (th *TestHelper) CreateTestBlog() models.Blog {
	return models.Blog{
		ID:          "test-blog-123",
		Title:       "Test Blog Title",
		Content:     "This is a test blog content for testing purposes.",
		AuthorID:    "test-user-123",
		AuthorName:  "test@example.com",
		Tags:        []string{"test", "go", "blog"},
		ViewCount:   0,
		Likes:       0,
		Dislikes:    0,
		Comments:    []models.Comment{},
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestBlogs creates multiple sample blogs for testing
func (th *TestHelper) CreateTestBlogs(count int) []models.Blog {
	blogs := make([]models.Blog, count)
	for i := 0; i < count; i++ {
		blogs[i] = models.Blog{
			ID:          fmt.Sprintf("test-blog-%d", i+1),
			Title:       fmt.Sprintf("Test Blog %d", i+1),
			Content:     fmt.Sprintf("This is test blog content %d", i+1),
			AuthorID:    fmt.Sprintf("test-user-%d", i+1),
			AuthorName:  fmt.Sprintf("user%d@example.com", i+1),
			Tags:        []string{"test", "go"},
			ViewCount:   i * 10,
			Likes:       i * 5,
			Dislikes:    i,
			Comments:    []models.Comment{},
			IsPublished: true,
			CreatedAt:   time.Now().AddDate(0, 0, -i),
			UpdatedAt:   time.Now().AddDate(0, 0, -i),
		}
	}
	return blogs
}

// CreateTestComment creates a sample comment for testing
func (th *TestHelper) CreateTestComment() models.Comment {
	return models.Comment{
		ID:         "test-comment-123",
		BlogID:     "test-blog-123",
		AuthorID:   "test-user-456",
		AuthorName: "commenter@example.com",
		Content:    "This is a test comment for testing purposes.",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// CreateTestComments creates multiple sample comments for testing
func (th *TestHelper) CreateTestComments(count int) []models.Comment {
	comments := make([]models.Comment, count)
	for i := 0; i < count; i++ {
		comments[i] = models.Comment{
			ID:         fmt.Sprintf("test-comment-%d", i+1),
			BlogID:     "test-blog-123",
			AuthorID:   fmt.Sprintf("test-user-%d", i+1),
			AuthorName: fmt.Sprintf("commenter%d@example.com", i+1),
			Content:    fmt.Sprintf("This is test comment %d", i+1),
			CreatedAt:  time.Now().AddDate(0, 0, -i),
			UpdatedAt:  time.Now().AddDate(0, 0, -i),
		}
	}
	return comments
}

// CreateTestUser creates a sample user for testing
func (th *TestHelper) CreateTestUser() models.User {
	return models.User{
		ID:       "test-user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "user",
		Verified: true,
		Bio:      "Test user bio",
		Picture:  "test-picture.jpg",
		Contact:  "test-contact",
	}
}

// GetTestHelper returns a new TestHelper instance
func GetTestHelper() *TestHelper {
	return &TestHelper{}
}
