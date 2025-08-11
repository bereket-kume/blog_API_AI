package repositories

import (
	"blog-api/Domain/models"
	"blog-api/Infrastructure/db_models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoRepo struct {
	collection *mongo.Collection
}

func NewUserMongoRepo(col *mongo.Collection) *userMongoRepo {
	return &userMongoRepo{collection: col}
}

func (ur *userMongoRepo) Insert(user *models.User) error {
	user.Verified = true
	db_user := db_models.FromDomainUser(user)
	_, err := ur.collection.InsertOne(context.TODO(), db_user)
	return err
}

func (ur *userMongoRepo) FindByEmail(email string) (*models.User, error) {
	filter := bson.M{"email": email}
	var user db_models.UserModel
	err := ur.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return db_models.ToDomainUser(&user), err
}

func (ur *userMongoRepo) UpdatePass(email string, passwordHash string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password_hash": passwordHash}}
	_, err := ur.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (ur *userMongoRepo) UpdateRole(email string, role string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"role": role}}
	_, err := ur.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (ur *userMongoRepo) Delete(email string) error {
	_, err := ur.collection.DeleteOne(context.TODO(), bson.M{"email": email})
	return err
}

func (ur *userMongoRepo) Verify(email string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"verified": true}}
	_, err := ur.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (ur *userMongoRepo) CountUsers() (int64, error) {
	count, err := ur.collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetUserByID retrieves a user by their ID
func (ur *userMongoRepo) GetUserByID(ctx context.Context, id string) (models.User, error) {
	// Convert string ID to primitive.ObjectID for MongoDB
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("invalid user ID")
	}

	var user db_models.UserModel
	err = ur.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	domainUser := db_models.ToDomainUser(&user)
	return *domainUser, err
}

// UpdateUserProfile updates a user's profile information
func (ur *userMongoRepo) UpdateUserProfile(ctx context.Context, id string, updated models.User) (models.User, error) {
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

	_, err = ur.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return models.User{}, err
	}

	return ur.GetUserByID(ctx, id)
}
