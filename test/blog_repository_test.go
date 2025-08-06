package test

import (
	domain "blog-api/Domain/models"
	"blog-api/Infrastructure/repositories"
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositoriesTestSuite struct {
	suite.Suite
	Client         *mongo.Client
	Database       *mongo.Database
	BlogCollection *mongo.Collection
	Repository     *repositories.BlogRepository
}

func (suite *RepositoriesTestSuite) SetupTest() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	suite.Require().NoError(err)

	database := client.Database("test_db")
	collection := database.Collection("blogs")

	suite.Client = client
	suite.Database = database
	suite.BlogCollection = collection
	suite.Repository = repositories.NewBlogRepository(collection)

	// Clean up the collection before each test
	_, err = collection.DeleteMany(context.Background(), bson.M{})
	suite.Require().NoError(err)
}

func (suite *RepositoriesTestSuite) TearDownTest() {
	err := suite.Client.Disconnect(context.Background())
	suite.Require().NoError(err)
}

func (suite *RepositoriesTestSuite) TestCreateBlog() {
	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog",
		Tags:    []string{"test", "blog"},
	}

	createdBlog, err := suite.Repository.CreateBlog(blog)
	suite.Require().NoError(err)
	suite.NotEmpty(createdBlog.ID)
	suite.Equal(blog.Title, createdBlog.Title)
	suite.Equal(blog.Content, createdBlog.Content)
}

func (suite *RepositoriesTestSuite) TestGetPaginatedBlogs() {
	for i := 1; i <= 15; i++ {
		blog := domain.Blog{
			Title:   "Blog " + strconv.Itoa(i),
			Content: "Content " + strconv.Itoa(i),
		}
		_, err := suite.Repository.CreateBlog(blog)
		suite.Require().NoError(err)
	}

	blogs, err := suite.Repository.GetPaginatedBlogs(1, 10)
	suite.Require().NoError(err)
	suite.Len(blogs, 10)
}

func (suite *RepositoriesTestSuite) TestGetBlogByID() {
	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog",
	}
	createdBlog, err := suite.Repository.CreateBlog(blog)
	suite.Require().NoError(err)

	fetchedBlog, err := suite.Repository.GetBlogByID(createdBlog.ID.Hex())
	suite.Require().NoError(err)
	suite.Equal(createdBlog.ID, fetchedBlog.ID)
}

func (suite *RepositoriesTestSuite) TestUpdateBlog() {
	blog := domain.Blog{
		Title:   "Old Title",
		Content: "Old Content",
	}
	createdBlog, err := suite.Repository.CreateBlog(blog)
	suite.Require().NoError(err)

	createdBlog.Title = "New Title"
	createdBlog.Content = "New Content"
	updatedBlog, err := suite.Repository.UpdateBlog(createdBlog)
	suite.Require().NoError(err)
	suite.Equal("New Title", updatedBlog.Title)
	suite.Equal("New Content", updatedBlog.Content)
}

func (suite *RepositoriesTestSuite) TestDeleteBlog() {
	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog",
	}
	createdBlog, err := suite.Repository.CreateBlog(blog)
	suite.Require().NoError(err)

	err = suite.Repository.DeleteBlog(createdBlog.ID.Hex())
	suite.Require().NoError(err)

	_, err = suite.Repository.GetBlogByID(createdBlog.ID.Hex())
	suite.Error(err)
}

func (suite *RepositoriesTestSuite) TestSearchBlogs() {
	blog1 := domain.Blog{
		Title:   "Searchable Blog 1",
		Content: "Content 1",
	}
	blog2 := domain.Blog{
		Title:   "Searchable Blog 2",
		Content: "Content 2",
	}
	_, err := suite.Repository.CreateBlog(blog1)
	suite.Require().NoError(err)
	_, err = suite.Repository.CreateBlog(blog2)
	suite.Require().NoError(err)

	blogs, err := suite.Repository.SearchBlogs("Searchable")
	suite.Require().NoError(err)
	suite.Len(blogs, 2)
}

func (suite *RepositoriesTestSuite) TestFilterBlogs() {
	blog := domain.Blog{
		Title:     "Filtered Blog",
		Content:   "Content",
		Tags:      []string{"tag1", "tag2"},
		CreatedAt: time.Now(),
	}
	_, err := suite.Repository.CreateBlog(blog)
	suite.Require().NoError(err)

	blogs, err := suite.Repository.FilterBlogs([]string{"tag1"}, [2]string{"2023-01-01T00:00:00Z", "2023-12-31T23:59:59Z"}, "created_at")
	suite.Require().NoError(err)
	suite.Len(blogs, 1)
}

func (suite *RepositoriesTestSuite) TestIncrementViewCount() {
	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog",
	}
	createdBlog, err := suite.Repository.CreateBlog(blog)
	suite.Require().NoError(err)

	err = suite.Repository.IncrementViewCount(createdBlog.ID.Hex())
	suite.Require().NoError(err)

	fetchedBlog, err := suite.Repository.GetBlogByID(createdBlog.ID.Hex())
	suite.Require().NoError(err)
	suite.Equal(1, fetchedBlog.Popularity.Views)
}

func TestRepositoriesTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesTestSuite))
}
