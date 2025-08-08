package repositories

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	Database "blog-api/Infrastructure/database"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userRepository struct{}

func NewUserRepository() interfaces.UserRepository {
	return &userRepository{}
}

func (r *userRepository) UpdateUserProfile(ctx context.Context, id string, updated models.User) (models.User, error) {
	collection := Database.GetCollection("users")

	// Convert string ID to primitive.ObjectID for MongoDB
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"bio":     updated.Bio,
			"picture": updated.Picture,
			"contact": updated.Contact,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return models.User{}, err
	}

	return r.GetUserByID(ctx, id)
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (models.User, error) {
	collection := Database.GetCollection("users")

	// Convert string ID to primitive.ObjectID for MongoDB
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("invalid user ID")
	}

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	// Convert the ID back to string for domain model
	user.ID = objectID.Hex()
	return user, err
}

func (r *userRepository) CountUsers() (int64, error) {
	collection := Database.GetCollection("users")
	return collection.CountDocuments(context.TODO(), bson.M{})
}

// Implement remaining methods from interfaces.UserRepository
func (r *userRepository) Insert(user *models.User) error {
	collection := Database.GetCollection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	return err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	collection := Database.GetCollection("users")
	filter := bson.M{"email": email}
	var user models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *userRepository) UpdatePass(email string, passwordHash string) error {
	collection := Database.GetCollection("users")
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": passwordHash}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *userRepository) UpdateRole(email string, role string) error {
	collection := Database.GetCollection("users")
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"role": role}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *userRepository) Delete(email string) error {
	collection := Database.GetCollection("users")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"email": email})
	return err
}

func (r *userRepository) Verify(email string) error {
	collection := Database.GetCollection("users")
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"verified": true}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}
