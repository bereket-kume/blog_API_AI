package repositories

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	Database "blog-api/Infrastructure/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// repositories/user_repository.go

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(col *mongo.Collection) interfaces.UserRepository {
	return &userRepository{collection: col}
}

func (r *userRepository) UpdateUserProfile(ctx context.Context, id primitive.ObjectID, updated models.User) (models.User, error) {
	collection := Database.GetUserCollection()

	update := bson.M{
		"$set": bson.M{
			"bio":     updated.Bio,
			"picture": updated.Picture,
			"contact": updated.Contact,
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return models.User{}, err
	}

	return r.GetUserByID(ctx, id)
}

func (r *userRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	collection := Database.GetUserCollection()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return user, err
}
