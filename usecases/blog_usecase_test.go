package usecases

import (
	"blog-api/Domain/models"
	"blog-api/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBlogUseCase_CreateBlog(t *testing.T) {
	tests := []struct {
		name        string
		blog        models.Blog
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
	}{
		{
			name: "Success - Create blog",
			blog: models.Blog{
				Title:       "Test Blog",
				Content:     "Test Content",
				AuthorID:    "user123",
				AuthorName:  "test@example.com",
				Tags:        []string{"test", "go"},
				IsPublished: true,
			},
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
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
				mockRepo.On("CreateBlog", mock.AnythingOfType("models.Blog")).Return(expectedBlog, nil)
			},
			expectError: false,
		},
		{
			name: "Error - Repository error",
			blog: models.Blog{
				Title:       "Test Blog",
				Content:     "Test Content",
				AuthorID:    "user123",
				AuthorName:  "test@example.com",
				Tags:        []string{"test", "go"},
				IsPublished: true,
			},
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("CreateBlog", mock.AnythingOfType("models.Blog")).Return(models.Blog{}, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			result, err := useCase.CreateBlog(tt.blog)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result.ID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.ID)
				assert.Equal(t, tt.blog.Title, result.Title)
				assert.Equal(t, tt.blog.Content, result.Content)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_GetPaginatedBlogs(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		limit       int
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
		expectedLen int
	}{
		{
			name:  "Success - Get paginated blogs",
			page:  1,
			limit: 10,
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				blogs := []models.Blog{
					{ID: "blog1", Title: "Blog 1", Content: "Content 1"},
					{ID: "blog2", Title: "Blog 2", Content: "Content 2"},
				}
				mockRepo.On("GetPaginatedBlogs", 1, 10).Return(blogs, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:  "Error - Repository error",
			page:  1,
			limit: 10,
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("GetPaginatedBlogs", 1, 10).Return(nil, errors.New("database error"))
			},
			expectError: true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			result, err := useCase.GetPaginatedBlogs(tt.page, tt.limit)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_GetBlogByID(t *testing.T) {
	tests := []struct {
		name        string
		blogID      string
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
	}{
		{
			name:   "Success - Get blog by ID",
			blogID: "blog123",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				blog := models.Blog{
					ID:      "blog123",
					Title:   "Test Blog",
					Content: "Test Content",
				}
				mockRepo.On("GetBlogByID", "blog123").Return(blog, nil)
			},
			expectError: false,
		},
		{
			name:   "Error - Blog not found",
			blogID: "nonexistent",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("GetBlogByID", "nonexistent").Return(models.Blog{}, errors.New("blog not found"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			result, err := useCase.GetBlogByID(tt.blogID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result.ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.blogID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_UpdateBlog(t *testing.T) {
	tests := []struct {
		name        string
		blog        models.Blog
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
	}{
		{
			name: "Success - Update blog",
			blog: models.Blog{
				ID:      "blog123",
				Title:   "Updated Blog",
				Content: "Updated Content",
			},
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				updatedBlog := models.Blog{
					ID:        "blog123",
					Title:     "Updated Blog",
					Content:   "Updated Content",
					UpdatedAt: time.Now(),
				}
				mockRepo.On("UpdateBlog", mock.AnythingOfType("models.Blog")).Return(updatedBlog, nil)
			},
			expectError: false,
		},
		{
			name: "Error - Repository error",
			blog: models.Blog{
				ID:      "blog123",
				Title:   "Updated Blog",
				Content: "Updated Content",
			},
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("UpdateBlog", mock.AnythingOfType("models.Blog")).Return(models.Blog{}, errors.New("update failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			result, err := useCase.UpdateBlog(tt.blog)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result.ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.blog.ID, result.ID)
				assert.Equal(t, tt.blog.Title, result.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_DeleteBlog(t *testing.T) {
	tests := []struct {
		name        string
		blogID      string
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
	}{
		{
			name:   "Success - Delete blog",
			blogID: "blog123",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("DeleteBlog", "blog123").Return(nil)
			},
			expectError: false,
		},
		{
			name:   "Error - Repository error",
			blogID: "nonexistent",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("DeleteBlog", "nonexistent").Return(errors.New("delete failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			err := useCase.DeleteBlog(tt.blogID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_SearchBlogs(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
		expectedLen int
	}{
		{
			name:  "Success - Search blogs",
			query: "test",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				blogs := []models.Blog{
					{ID: "blog1", Title: "Test Blog 1"},
					{ID: "blog2", Title: "Test Blog 2"},
				}
				mockRepo.On("SearchBlogs", "test").Return(blogs, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:  "Error - Repository error",
			query: "test",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("SearchBlogs", "test").Return(nil, errors.New("search failed"))
			},
			expectError: true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			result, err := useCase.SearchBlogs(tt.query)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_FilterBlogs(t *testing.T) {
	tests := []struct {
		name        string
		tags        []string
		dateRange   [2]string
		sortBy      string
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
		expectedLen int
	}{
		{
			name:      "Success - Filter blogs",
			tags:      []string{"go", "test"},
			dateRange: [2]string{"2023-01-01", "2023-12-31"},
			sortBy:    "created_at",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				blogs := []models.Blog{
					{ID: "blog1", Title: "Blog 1", Tags: []string{"go"}},
					{ID: "blog2", Title: "Blog 2", Tags: []string{"test"}},
				}
				mockRepo.On("FilterBlogs", []string{"go", "test"}, [2]string{"2023-01-01", "2023-12-31"}, "created_at").Return(blogs, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:      "Error - Repository error",
			tags:      []string{"go"},
			dateRange: [2]string{},
			sortBy:    "title",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("FilterBlogs", []string{"go"}, [2]string{"", ""}, "title").Return(nil, errors.New("filter failed"))
			},
			expectError: true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			result, err := useCase.FilterBlogs(tt.tags, tt.dateRange, tt.sortBy)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_IncrementViewCount(t *testing.T) {
	tests := []struct {
		name        string
		blogID      string
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
	}{
		{
			name:   "Success - Increment view count",
			blogID: "blog123",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("IncrementViewCount", "blog123").Return(nil)
			},
			expectError: false,
		},
		{
			name:   "Error - Repository error",
			blogID: "nonexistent",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("IncrementViewCount", "nonexistent").Return(errors.New("increment failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			err := useCase.IncrementViewCount(tt.blogID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_UpdateLikes(t *testing.T) {
	tests := []struct {
		name        string
		blogID      string
		increment   bool
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
	}{
		{
			name:      "Success - Increment likes",
			blogID:    "blog123",
			increment: true,
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("UpdateLikes", "blog123", true).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "Success - Decrement likes",
			blogID:    "blog123",
			increment: false,
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("UpdateLikes", "blog123", false).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "Error - Repository error",
			blogID:    "nonexistent",
			increment: true,
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("UpdateLikes", "nonexistent", true).Return(errors.New("update failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			err := useCase.UpdateLikes(tt.blogID, tt.increment)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_UpdateDislikes(t *testing.T) {
	tests := []struct {
		name        string
		blogID      string
		increment   bool
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
	}{
		{
			name:      "Success - Increment dislikes",
			blogID:    "blog123",
			increment: true,
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("UpdateDislikes", "blog123", true).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "Success - Decrement dislikes",
			blogID:    "blog123",
			increment: false,
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("UpdateDislikes", "blog123", false).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "Error - Repository error",
			blogID:    "nonexistent",
			increment: true,
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("UpdateDislikes", "nonexistent", true).Return(errors.New("update failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			err := useCase.UpdateDislikes(tt.blogID, tt.increment)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_AddComment(t *testing.T) {
	tests := []struct {
		name        string
		blogID      string
		comment     models.Comment
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
	}{
		{
			name:   "Success - Add comment",
			blogID: "blog123",
			comment: models.Comment{
				AuthorID:   "user123",
				AuthorName: "test@example.com",
				Content:    "Great blog post!",
			},
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				expectedComment := models.Comment{
					ID:         "comment123",
					BlogID:     "blog123",
					AuthorID:   "user123",
					AuthorName: "test@example.com",
					Content:    "Great blog post!",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				mockRepo.On("AddComment", "blog123", mock.AnythingOfType("models.Comment")).Return(expectedComment, nil)
			},
			expectError: false,
		},
		{
			name:   "Error - Repository error",
			blogID: "blog123",
			comment: models.Comment{
				AuthorID:   "user123",
				AuthorName: "test@example.com",
				Content:    "Great blog post!",
			},
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("AddComment", "blog123", mock.AnythingOfType("models.Comment")).Return(models.Comment{}, errors.New("add comment failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			result, err := useCase.AddComment(tt.blogID, tt.comment)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result.ID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.ID)
				assert.Equal(t, tt.blogID, result.BlogID)
				assert.Equal(t, tt.comment.Content, result.Content)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBlogUseCase_GetComments(t *testing.T) {
	tests := []struct {
		name        string
		blogID      string
		setupMock   func(*mocks.BlogRepositoryMock)
		expectError bool
		expectedLen int
	}{
		{
			name:   "Success - Get comments",
			blogID: "blog123",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				comments := []models.Comment{
					{ID: "comment1", BlogID: "blog123", Content: "Comment 1"},
					{ID: "comment2", BlogID: "blog123", Content: "Comment 2"},
				}
				mockRepo.On("GetComments", "blog123").Return(comments, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:   "Error - Repository error",
			blogID: "blog123",
			setupMock: func(mockRepo *mocks.BlogRepositoryMock) {
				mockRepo.On("GetComments", "blog123").Return(nil, errors.New("get comments failed"))
			},
			expectError: true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.BlogRepositoryMock{}
			tt.setupMock(mockRepo)

			useCase := NewBlogUseCase(mockRepo)
			result, err := useCase.GetComments(tt.blogID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
