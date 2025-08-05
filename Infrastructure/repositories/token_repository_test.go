package repositories

import (
	"blog-api/Domain/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestTokenMongoRepo(t *testing.T) {
	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	m.Run("Insert", func(mt *mtest.T) {
		repo := NewTokenMongoRepo(mt.Coll)
		token := models.Token{
			ID:        "507f1f77bcf86cd799439011",
			UserID:    "testuser",
			Token:     "token-string",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now(),
			IP:        "1234.5.6",
			Device:    "Chrome on window",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := repo.CreateToken(token)
		assert.NoError(mt, err)
	})
	m.Run("FindByID", func(mt *mtest.T) {
		repo := NewTokenMongoRepo(mt.Coll)
		expected := models.Token{
			ID:        "507f1f77bcf86cd799439011",
			UserID:    "testuser",
			Token:     "token-string",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now(),
			IP:        "1234.5.6",
			Device:    "Chrome on window",
		}

		// Create a mock response document
		first := mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
			{"_id", expected.ID},
			{"user_id", expected.UserID},
			{"token_hash", expected.Token},
			{"created_at", expected.CreatedAt},
			{"expired_at", expected.ExpiresAt},
			{"ip", expected.IP},
			{"device", expected.Device},
		})
		mt.AddMockResponses(first)

		foundUser, err := repo.GetToken(expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected.UserID, foundUser.UserID)
	})
	m.Run("Update", func(mt *mtest.T) {
		repo := NewTokenMongoRepo(mt.Coll)
		// Simulate MongoDB's UpdateOne success response

		token := models.Token{
			ID:        "507f1f77bcf86cd799439011",
			UserID:    "testuser",
			Token:     "token-string",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now(),
			IP:        "1234.5.6",
			Device:    "Chrome on window",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse()) // mock successful update

		err := repo.Update(token)

		assert.NoError(t, err)

	})
	m.Run("Delete", func(mt *mtest.T) {
		repo := NewUserMongoRepo(mt.Coll)
		mt.AddMockResponses(mtest.CreateSuccessResponse()) // Simulates successful deletion

		err := repo.Delete("1")
		assert.NoError(t, err)

	})
}

// package repositories

// import (
// 	"blog-api/Domain/models"
// 	"context"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	mongodriver "go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func setupTestDB(t *testing.T) (*mongodriver.Client, *mongodriver.Collection, func()) {
// 	// Connect to MongoDB running in Docker
// 	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
// 	client, err := mongodriver.Connect(context.TODO(), clientOpts)
// 	assert.NoError(t, err)

// 	db := client.Database("testdb")
// 	collection := db.Collection("users_test")

// 	// Cleanup function
// 	cleanup := func() {
// 		_ = collection.Drop(context.TODO()) // Drop the collection
// 		_ = client.Disconnect(context.TODO())
// 	}

// 	return client, collection, cleanup
// }
// func TestMongoRepo(t *testing.T) {
// 	_, collection, cleanup := setupTestDB(t)
// 	defer cleanup() // Runs after all tests in this function

// 	repo := NewUserMongoRepo(collection)

// 	user := models.User{
// 		ID:       "507f1f77bcf86cd799439011",
// 		Username: "testuser",
// 		Email:    "test1@example.com",
// 		Password: "pass",
// 		Role:     "user",
// 		Verified: false,
// 	}

// 	err := repo.Insert(user)
// 	assert.NoError(t, err)

// 	foundUser, err := repo.FindByEmail("test1@example.com")
// 	assert.NoError(t, err)
// 	assert.Equal(t, user.Email, foundUser.Email)

// 	err = repo.UpdatePass("test1@example.com", "changed_pass")
// 	assert.NoError(t, err)

// 	updatedUser, err := repo.FindByEmail("test1@example.com")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "changed_pass", updatedUser.Password)

// 	err = repo.UpdateRole("test1@example.com", "admin")
// 	assert.NoError(t, err)

// 	updatedUser, err = repo.FindByEmail("test1@example.com")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "admin", updatedUser.Role)

// 	err = repo.Verify("test1@example.com")
// 	assert.NoError(t, err)

// 	verifiedUser, err := repo.FindByEmail("test1@example.com")
// 	assert.NoError(t, err)
// 	assert.Equal(t, true, verifiedUser.Verified)

// }
