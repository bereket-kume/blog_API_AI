package repositories

import (
	"blog-api/Domain/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// Example assumes your Insert method uses collection.InsertOne

func TestUserMongoRepo(t *testing.T) {
	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	m.Run("Insert", func(mt *mtest.T) {
		repo := NewUserMongoRepo(mt.Coll)
		user := models.User{
			ID:       "507f1f77bcf86cd799439011",
			Username: "testuser",
			Email:    "test1@example.com",
			Password: "pass",
			Role:     "user",
			Verified: false,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := repo.Insert(user)
		assert.NoError(mt, err)
	})
	m.Run("FindByEmail", func(mt *mtest.T) {
		repo := NewUserMongoRepo(mt.Coll)
		expected := models.User{
			ID:       "507f1f77bcf86cd799439011",
			Username: "testuser",
			Email:    "test1@example.com",
			Password: "pass",
			Role:     "user",
			Verified: false,
		}

		// Create a mock response document
		first := mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
			{"_id", expected.ID},
			{"username", expected.Username},
			{"email", expected.Email},
			{"password", expected.Password},
			{"role", expected.Role},
			{"verified", expected.Verified},
		})
		mt.AddMockResponses(first)

		foundUser, err := repo.FindByEmail(expected.Email)
		assert.NoError(t, err)
		assert.Equal(t, expected.Email, foundUser.Email)
	})
	m.Run("UpdatePass Success", func(mt *mtest.T) {
		repo := NewUserMongoRepo(mt.Coll)
		// Simulate MongoDB's UpdateOne success response
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: 1},
			bson.E{Key: "nModified", Value: 1},
			bson.E{Key: "ok", Value: 1},
		))

		err := repo.UpdatePass("test@example.com", "newHashedPassword")
		assert.NoError(t, err)
	})
	m.Run("UpdateRole Success", func(mt *mtest.T) {
		repo := NewUserMongoRepo(mt.Coll)
		// Simulate MongoDB's UpdateOne success response
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: 1},
			bson.E{Key: "nModified", Value: 1},
			bson.E{Key: "ok", Value: 1},
		))

		err := repo.UpdateRole("test@example.com", "admin")
		assert.NoError(t, err)
	})

	m.Run("Verify Success", func(mt *mtest.T) {
		repo := NewUserMongoRepo(mt.Coll)
		// Simulate MongoDB's UpdateOne success response
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: 1},
			bson.E{Key: "nModified", Value: 1},
			bson.E{Key: "ok", Value: 1},
		))

		err := repo.Verify("test@example.com")
		assert.NoError(t, err)
	})
	m.Run("CountUsers Success", func(mt *mtest.T) {
		repo := NewUserMongoRepo(mt.Coll)
		// Simulate a successful count response (e.g. 5 users)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch, bson.D{
			{"n", int32(5)},
		}))

		count, err := repo.CountUsers()
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
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
