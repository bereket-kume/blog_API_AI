package mocks

import (
	"blog-api/Domain/models"

	"github.com/stretchr/testify/mock"
)

type BlogUseCaseMock struct {
	mock.Mock
}

func (m *BlogUseCaseMock) CreateBlog(blog models.Blog) (models.Blog, error) {
	args := m.Called(blog)
	return args.Get(0).(models.Blog), args.Error(1)
}

func (m *BlogUseCaseMock) GetPaginatedBlogs(page, limit int) ([]models.Blog, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Blog), args.Error(1)
}

func (m *BlogUseCaseMock) GetBlogByID(blogID string) (models.Blog, error) {
	args := m.Called(blogID)
	return args.Get(0).(models.Blog), args.Error(1)
}

func (m *BlogUseCaseMock) UpdateBlog(blog models.Blog) (models.Blog, error) {
	args := m.Called(blog)
	return args.Get(0).(models.Blog), args.Error(1)
}

func (m *BlogUseCaseMock) DeleteBlog(blogID string) error {
	args := m.Called(blogID)
	return args.Error(0)
}

func (m *BlogUseCaseMock) SearchBlogs(query string) ([]models.Blog, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Blog), args.Error(1)
}

func (m *BlogUseCaseMock) FilterBlogs(tags []string, dateRange [2]string, sortBy string) ([]models.Blog, error) {
	args := m.Called(tags, dateRange, sortBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Blog), args.Error(1)
}

func (m *BlogUseCaseMock) IncrementViewCount(blogID string) error {
	args := m.Called(blogID)
	return args.Error(0)
}

func (m *BlogUseCaseMock) UpdateLikes(blogID string, increment bool) error {
	args := m.Called(blogID, increment)
	return args.Error(0)
}

func (m *BlogUseCaseMock) UpdateDislikes(blogID string, increment bool) error {
	args := m.Called(blogID, increment)
	return args.Error(0)
}

func (m *BlogUseCaseMock) AddComment(blogID string, comment models.Comment) (models.Comment, error) {
	args := m.Called(blogID, comment)
	return args.Get(0).(models.Comment), args.Error(1)
}

func (m *BlogUseCaseMock) GetComments(blogID string) ([]models.Comment, error) {
	args := m.Called(blogID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Comment), args.Error(1)
}
