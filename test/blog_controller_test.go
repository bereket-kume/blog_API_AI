package test

import (
	domain "blog-api/Domain/models"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockBlogUseCase struct {
	mock.Mock
}

func (m *MockBlogUseCase) CreateBlog(blog domain.Blog) (domain.Blog, error) {
	args := m.Called(blog)
	return args.Get(0).(domain.Blog), args.Error(1)
}

func (m *MockBlogUseCase) GetPaginatedBlogs(page, limit int) ([]domain.Blog, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]domain.Blog), args.Error(1)
}

func (m *MockBlogUseCase) GetBlogByID(blogID string) (domain.Blog, error) {
	args := m.Called(blogID)
	return args.Get(0).(domain.Blog), args.Error(1)
}

func (m *MockBlogUseCase) UpdateBlog(blog domain.Blog) (domain.Blog, error) {
	args := m.Called(blog)
	return args.Get(0).(domain.Blog), args.Error(1)
}

func (m *MockBlogUseCase) DeleteBlog(blogID string) error {
	args := m.Called(blogID)
	return args.Error(0)
}

type ControllerTestSuite struct {
	suite.Suite
	BlogUseCaseMock *MockBlogUseCase
}

func (suite *ControllerTestSuite) SetupTest() {
	suite.BlogUseCaseMock = new(MockBlogUseCase)
}

func (suite *ControllerTestSuite) TestCreateBlog() {
	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog",
	}
	expectedBlog := domain.Blog{
		ID:      primitive.NewObjectID(),
		Title:   "Test Blog",
		Content: "This is a test blog",
	}
	suite.BlogUseCaseMock.On("CreateBlog", blog).Return(expectedBlog, nil)

	createdBlog, err := suite.BlogUseCaseMock.CreateBlog(blog)
	suite.NoError(err)
	suite.Equal(expectedBlog, createdBlog)
	suite.BlogUseCaseMock.AssertExpectations(suite.T())
}

func (suite *ControllerTestSuite) TestGetPaginatedBlogs() {
	expectedBlogs := []domain.Blog{
		{ID: primitive.NewObjectID(), Title: "Blog 1", Content: "Content 1"},
		{ID: primitive.NewObjectID(), Title: "Blog 2", Content: "Content 2"},
	}
	suite.BlogUseCaseMock.On("GetPaginatedBlogs", 1, 10).Return(expectedBlogs, nil)

	blogs, err := suite.BlogUseCaseMock.GetPaginatedBlogs(1, 10)
	suite.NoError(err)
	suite.Equal(len(expectedBlogs), len(blogs))
	suite.Equal(expectedBlogs[0].Title, blogs[0].Title)
	suite.BlogUseCaseMock.AssertExpectations(suite.T())
}

func (suite *ControllerTestSuite) TestGetBlogByID() {
	blogID := primitive.NewObjectID().Hex()
	expectedBlog := domain.Blog{ID: primitive.NewObjectID(), Title: "Test Blog", Content: "Test Content"}
	suite.BlogUseCaseMock.On("GetBlogByID", blogID).Return(expectedBlog, nil)

	blog, err := suite.BlogUseCaseMock.GetBlogByID(blogID)
	suite.NoError(err)
	suite.Equal(expectedBlog, blog)
	suite.BlogUseCaseMock.AssertExpectations(suite.T())
}

func (suite *ControllerTestSuite) TestUpdateBlog() {
	blog := domain.Blog{
		ID:      primitive.NewObjectID(),
		Title:   "Updated Blog",
		Content: "Updated Content",
	}
	suite.BlogUseCaseMock.On("UpdateBlog", blog).Return(blog, nil)

	updatedBlog, err := suite.BlogUseCaseMock.UpdateBlog(blog)
	suite.NoError(err)
	suite.Equal(blog.Title, updatedBlog.Title)
	suite.Equal(blog.Content, updatedBlog.Content)
	suite.BlogUseCaseMock.AssertExpectations(suite.T())
}

func (suite *ControllerTestSuite) TestDeleteBlog() {
	blogID := primitive.NewObjectID().Hex()
	suite.BlogUseCaseMock.On("DeleteBlog", blogID).Return(nil)

	err := suite.BlogUseCaseMock.DeleteBlog(blogID)
	suite.NoError(err)
	suite.BlogUseCaseMock.AssertExpectations(suite.T())
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}
