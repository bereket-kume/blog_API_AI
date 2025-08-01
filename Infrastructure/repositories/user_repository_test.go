package repositories

import (
	"blog-api/Domain/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongodriver.Client, *mongodriver.Collection, func()) {
	// Connect to MongoDB running in Docker
	clientOpts := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongodriver.Connect(context.TODO(), clientOpts)
	assert.NoError(t, err)

	db := client.Database("testdb")
	collection := db.Collection("users_test")

	// Cleanup function
	cleanup := func() {
		_ = collection.Drop(context.TODO()) // Drop the collection
		_ = client.Disconnect(context.TODO())
	}

	return client, collection, cleanup
}
func TestMongoRepo(t *testing.T) {
	_, collection, cleanup := setupTestDB(t)
	defer cleanup() // Runs after all tests in this function

	repo := NewUserMongoRepo(collection)

	user := models.User{
		ID:       "507f1f77bcf86cd799439011",
		Username: "testuser",
		Email:    "test1@example.com",
		Password: "pass",
		Role:     "user",
		Verified: false,
	}

	err := repo.Insert(user)
	assert.NoError(t, err)

	foundUser, err := repo.FindByEmail("test1@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, foundUser.Email)

	err = repo.UpdatePass("test1@example.com", "changed_pass")
	assert.NoError(t, err)

	updatedUser, err := repo.FindByEmail("test1@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "changed_pass", updatedUser.Password)

	err = repo.UpdateRole("test1@example.com", "admin")
	assert.NoError(t, err)

	updatedUser, err = repo.FindByEmail("test1@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "admin", updatedUser.Role)

	err = repo.Verify("test1@example.com")
	assert.NoError(t, err)

	verifiedUser, err := repo.FindByEmail("test1@example.com")
	assert.NoError(t, err)
	assert.Equal(t, true, verifiedUser.Verified)

}
