package repositories

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	Database "blog-api/Infrastructure/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userRepository struct{}

func NewUserRepository() interfaces.UserRepository {
	return &userRepository{}
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

func (r *userRepository) CountUsers() (int64, error) {
	collection := Database.GetUserCollection()
	return collection.CountDocuments(context.TODO(), bson.M{})
}
