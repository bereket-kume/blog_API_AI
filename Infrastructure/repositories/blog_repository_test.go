package repositories

import (
	"blog-api/Domain/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	repo       *blogMongoRepo
	ctx        context.Context
}

func (suite *BlogRepositoryTestSuite) SetupSuite() {
	// Connect to test database
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		suite.T().Fatalf("Failed to connect to MongoDB: %v", err)
	}
	suite.client = client

	// Create test database and collection
	suite.database = client.Database("blog_test_db")
	suite.collection = suite.database.Collection("blogs")
	suite.ctx = context.Background()

	// Initialize repository
	suite.repo = NewBlogMongoRepo(suite.collection)
}

func (suite *BlogRepositoryTestSuite) TearDownSuite() {
	// Drop test database
	if suite.database != nil {
		suite.database.Drop(suite.ctx)
	}

	// Close connection
	if suite.client != nil {
		suite.client.Disconnect(suite.ctx)
	}
}

func (suite *BlogRepositoryTestSuite) SetupTest() {
	// Clear collection before each test
	suite.collection.DeleteMany(suite.ctx, bson.M{})
}

func (suite *BlogRepositoryTestSuite) TearDownTest() {
	// Clear collection after each test
	suite.collection.DeleteMany(suite.ctx, bson.M{})
}

func (suite *BlogRepositoryTestSuite) TestCreateBlog() {
	// Test data
	blog := models.Blog{
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Tags:        []string{"test", "go"},
		Comments:    []models.Comment{}, // Initialize empty slice
		IsPublished: true,
	}

	// Create blog
	createdBlog, err := suite.repo.CreateBlog(blog)

	// Assertions
	suite.NoError(err)
	suite.NotEmpty(createdBlog.ID)
	suite.Equal(blog.Title, createdBlog.Title)
	suite.Equal(blog.Content, createdBlog.Content)
	suite.Equal(blog.AuthorID, createdBlog.AuthorID)
	suite.Equal(blog.AuthorName, createdBlog.AuthorName)
	suite.Equal(blog.Tags, createdBlog.Tags)
	suite.Equal(blog.IsPublished, createdBlog.IsPublished)
	suite.False(createdBlog.CreatedAt.IsZero())
	suite.False(createdBlog.UpdatedAt.IsZero())

	// Verify blog was saved in database
	var savedBlog models.Blog
	objectID, _ := primitive.ObjectIDFromHex(createdBlog.ID)
	err = suite.collection.FindOne(suite.ctx, bson.M{"_id": objectID}).Decode(&savedBlog)
	suite.NoError(err)
	suite.Equal(createdBlog.ID, savedBlog.ID)
}

func (suite *BlogRepositoryTestSuite) TestGetBlogByID() {
	// Create test blog
	blog := models.Blog{
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Tags:        []string{"test", "go"},
		Comments:    []models.Comment{}, // Initialize empty slice
		IsPublished: true,
	}
	createdBlog, _ := suite.repo.CreateBlog(blog)

	// Get blog by ID
	retrievedBlog, err := suite.repo.GetBlogByID(createdBlog.ID)

	// Assertions
	suite.NoError(err)
	suite.Equal(createdBlog.ID, retrievedBlog.ID)
	suite.Equal(createdBlog.Title, retrievedBlog.Title)
	suite.Equal(createdBlog.Content, retrievedBlog.Content)
}

func (suite *BlogRepositoryTestSuite) TestGetBlogByID_NotFound() {
	// Try to get non-existent blog
	nonExistentID := primitive.NewObjectID().Hex()
	_, err := suite.repo.GetBlogByID(nonExistentID)

	// Assertions
	suite.Error(err)
}

func (suite *BlogRepositoryTestSuite) TestGetBlogByID_InvalidID() {
	// Try to get blog with invalid ID
	_, err := suite.repo.GetBlogByID("invalid-id")

	// Assertions
	suite.Error(err)
}

func (suite *BlogRepositoryTestSuite) TestUpdateBlog() {
	// Create test blog
	blog := models.Blog{
		Title:       "Original Title",
		Content:     "Original Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Tags:        []string{"original"},
		Comments:    []models.Comment{}, // Initialize empty slice
		IsPublished: true,
	}
	createdBlog, _ := suite.repo.CreateBlog(blog)

	// Update blog
	createdBlog.Title = "Updated Title"
	createdBlog.Content = "Updated Content"
	createdBlog.Tags = []string{"updated", "go"}

	updatedBlog, err := suite.repo.UpdateBlog(createdBlog)

	// Assertions
	suite.NoError(err)
	suite.Equal("Updated Title", updatedBlog.Title)
	suite.Equal("Updated Content", updatedBlog.Content)
	suite.Equal([]string{"updated", "go"}, updatedBlog.Tags)
	suite.True(updatedBlog.UpdatedAt.After(createdBlog.UpdatedAt))

	// Verify in database
	retrievedBlog, err := suite.repo.GetBlogByID(createdBlog.ID)
	suite.NoError(err)
	suite.Equal("Updated Title", retrievedBlog.Title)
	suite.Equal("Updated Content", retrievedBlog.Content)
}

func (suite *BlogRepositoryTestSuite) TestDeleteBlog() {
	// Create test blog
	blog := models.Blog{
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Tags:        []string{"test"},
		IsPublished: true,
	}
	createdBlog, _ := suite.repo.CreateBlog(blog)

	// Delete blog
	err := suite.repo.DeleteBlog(createdBlog.ID)

	// Assertions
	suite.NoError(err)

	// Verify blog was deleted
	_, err = suite.repo.GetBlogByID(createdBlog.ID)
	suite.Error(err)
}

func (suite *BlogRepositoryTestSuite) TestGetPaginatedBlogs() {
	// Create multiple test blogs
	blogs := []models.Blog{
		{Title: "Blog 1", Content: "Content 1", AuthorID: "user1", AuthorName: "user1@test.com", IsPublished: true},
		{Title: "Blog 2", Content: "Content 2", AuthorID: "user2", AuthorName: "user2@test.com", IsPublished: true},
		{Title: "Blog 3", Content: "Content 3", AuthorID: "user3", AuthorName: "user3@test.com", IsPublished: true},
		{Title: "Draft Blog", Content: "Draft Content", AuthorID: "user4", AuthorName: "user4@test.com", IsPublished: false},
	}

	for _, blog := range blogs {
		suite.repo.CreateBlog(blog)
	}

	// Test pagination
	retrievedBlogs, err := suite.repo.GetPaginatedBlogs(1, 2)

	// Assertions
	suite.NoError(err)
	suite.Len(retrievedBlogs, 2)                   // Only published blogs should be returned
	suite.Equal("Blog 3", retrievedBlogs[0].Title) // Should be sorted by created_at desc
	suite.Equal("Blog 2", retrievedBlogs[1].Title)
}

func (suite *BlogRepositoryTestSuite) TestSearchBlogs() {
	// Create test blogs
	blogs := []models.Blog{
		{Title: "Go Programming", Content: "Learn Go programming", AuthorID: "user1", AuthorName: "user1@test.com", IsPublished: true},
		{Title: "Python Tutorial", Content: "Learn Python programming", AuthorID: "user2", AuthorName: "user2@test.com", IsPublished: true},
		{Title: "Web Development", Content: "Learn web development with Go", AuthorID: "user3", AuthorName: "user3@test.com", IsPublished: true},
		{Title: "Draft Blog", Content: "Draft content about Go", AuthorID: "user4", AuthorName: "user4@test.com", IsPublished: false},
	}

	for _, blog := range blogs {
		suite.repo.CreateBlog(blog)
	}

	// Search for "Go"
	results, err := suite.repo.SearchBlogs("Go")

	// Assertions
	suite.NoError(err)
	suite.Len(results, 2) // Should find "Go Programming" and "Web Development" (published only)

	// Verify search results
	titles := []string{results[0].Title, results[1].Title}
	suite.Contains(titles, "Go Programming")
	suite.Contains(titles, "Web Development")
}

func (suite *BlogRepositoryTestSuite) TestFilterBlogs() {
	// Create test blogs with different tags and dates
	now := time.Now()
	blogs := []models.Blog{
		{
			Title: "Go Blog", Content: "Go content", AuthorID: "user1", AuthorName: "user1@test.com",
			Tags: []string{"go", "programming"}, IsPublished: true, CreatedAt: now.AddDate(0, 0, -1),
		},
		{
			Title: "Python Blog", Content: "Python content", AuthorID: "user2", AuthorName: "user2@test.com",
			Tags: []string{"python", "programming"}, IsPublished: true, CreatedAt: now.AddDate(0, 0, -2),
		},
		{
			Title: "Web Blog", Content: "Web content", AuthorID: "user3", AuthorName: "user3@test.com",
			Tags: []string{"web", "frontend"}, IsPublished: true, CreatedAt: now.AddDate(0, 0, -3),
		},
	}

	for _, blog := range blogs {
		suite.repo.CreateBlog(blog)
	}

	// Test filtering by tags
	results, err := suite.repo.FilterBlogs([]string{"go"}, [2]string{}, "created_at")

	// Assertions
	suite.NoError(err)
	suite.Len(results, 1)
	suite.Equal("Go Blog", results[0].Title)
}

func (suite *BlogRepositoryTestSuite) TestIncrementViewCount() {
	// Create test blog
	blog := models.Blog{
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		ViewCount:   0,
		IsPublished: true,
	}
	createdBlog, _ := suite.repo.CreateBlog(blog)

	// Increment view count
	err := suite.repo.IncrementViewCount(createdBlog.ID)

	// Assertions
	suite.NoError(err)

	// Verify view count was incremented
	updatedBlog, _ := suite.repo.GetBlogByID(createdBlog.ID)
	suite.Equal(1, updatedBlog.ViewCount)
}

func (suite *BlogRepositoryTestSuite) TestUpdateLikes() {
	// Create test blog
	blog := models.Blog{
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Likes:       0,
		IsPublished: true,
	}
	createdBlog, _ := suite.repo.CreateBlog(blog)

	// Test increment likes
	err := suite.repo.UpdateLikes(createdBlog.ID, true)
	suite.NoError(err)

	updatedBlog, _ := suite.repo.GetBlogByID(createdBlog.ID)
	suite.Equal(1, updatedBlog.Likes)

	// Test decrement likes
	err = suite.repo.UpdateLikes(createdBlog.ID, false)
	suite.NoError(err)

	updatedBlog, _ = suite.repo.GetBlogByID(createdBlog.ID)
	suite.Equal(0, updatedBlog.Likes)
}

func (suite *BlogRepositoryTestSuite) TestUpdateDislikes() {
	// Create test blog
	blog := models.Blog{
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Dislikes:    0,
		IsPublished: true,
	}
	createdBlog, _ := suite.repo.CreateBlog(blog)

	// Test increment dislikes
	err := suite.repo.UpdateDislikes(createdBlog.ID, true)
	suite.NoError(err)

	updatedBlog, _ := suite.repo.GetBlogByID(createdBlog.ID)
	suite.Equal(1, updatedBlog.Dislikes)

	// Test decrement dislikes
	err = suite.repo.UpdateDislikes(createdBlog.ID, false)
	suite.NoError(err)

	updatedBlog, _ = suite.repo.GetBlogByID(createdBlog.ID)
	suite.Equal(0, updatedBlog.Dislikes)
}

func (suite *BlogRepositoryTestSuite) TestAddComment() {
	// Create test blog
	blog := models.Blog{
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Comments:    []models.Comment{}, // Initialize empty slice
		IsPublished: true,
	}
	createdBlog, _ := suite.repo.CreateBlog(blog)

	// Add comment
	comment := models.Comment{
		AuthorID:   "commenter123",
		AuthorName: "commenter@test.com",
		Content:    "Great blog post!",
	}

	createdComment, err := suite.repo.AddComment(createdBlog.ID, comment)

	// Assertions
	suite.NoError(err)
	suite.NotEmpty(createdComment.ID)
	suite.Equal(createdBlog.ID, createdComment.BlogID)
	suite.Equal(comment.AuthorID, createdComment.AuthorID)
	suite.Equal(comment.AuthorName, createdComment.AuthorName)
	suite.Equal(comment.Content, createdComment.Content)
	suite.False(createdComment.CreatedAt.IsZero())
	suite.False(createdComment.UpdatedAt.IsZero())

	// Verify comment was added to blog
	updatedBlog, _ := suite.repo.GetBlogByID(createdBlog.ID)
	suite.Len(updatedBlog.Comments, 1)
	suite.Equal(createdComment.ID, updatedBlog.Comments[0].ID)
}

func (suite *BlogRepositoryTestSuite) TestGetComments() {
	// Create test blog with comments
	blog := models.Blog{
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Comments:    []models.Comment{}, // Initialize empty slice
		IsPublished: true,
	}
	createdBlog, _ := suite.repo.CreateBlog(blog)

	// Add comments
	comments := []models.Comment{
		{AuthorID: "user1", AuthorName: "user1@test.com", Content: "Comment 1"},
		{AuthorID: "user2", AuthorName: "user2@test.com", Content: "Comment 2"},
	}

	for _, comment := range comments {
		suite.repo.AddComment(createdBlog.ID, comment)
	}

	// Get comments
	retrievedComments, err := suite.repo.GetComments(createdBlog.ID)

	// Assertions
	suite.NoError(err)
	suite.Len(retrievedComments, 2)
	if len(retrievedComments) >= 1 {
		suite.Equal("Comment 1", retrievedComments[0].Content)
	}
	if len(retrievedComments) >= 2 {
		suite.Equal("Comment 2", retrievedComments[1].Content)
	}
}

func (suite *BlogRepositoryTestSuite) TestGetComments_NoComments() {
	// Create test blog without comments
	blog := models.Blog{
		Title:       "Test Blog",
		Content:     "Test Content",
		AuthorID:    "user123",
		AuthorName:  "test@example.com",
		Comments:    []models.Comment{}, // Initialize empty slice
		IsPublished: true,
	}
	createdBlog, _ := suite.repo.CreateBlog(blog)

	// Get comments
	retrievedComments, err := suite.repo.GetComments(createdBlog.ID)

	// Assertions
	suite.NoError(err)
	suite.Len(retrievedComments, 0)
}

// Run the test suite
func TestBlogRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BlogRepositoryTestSuite))
}
