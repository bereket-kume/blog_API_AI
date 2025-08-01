package repositories

import (
	"blog-api/Domain/models"
	"blog-api/Infrastructure/db_models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoRepo struct {
	collection *mongo.Collection
}

func NewUserMongoRepo(col *mongo.Collection) *userMongoRepo {
	return &userMongoRepo{collection: col}
}

func (ur *userMongoRepo) Insert(user models.User) error {
	db_user := db_models.FromDomainUser(user)
	_, err := ur.collection.InsertOne(context.TODO(), db_user)
	return err
}

func (ur *userMongoRepo) FindByEmail(email string) (models.User, error) {
	filter := bson.M{"email": email}
	var user db_models.UserModel
	err := ur.collection.FindOne(context.TODO(), filter).Decode(&user)
	return db_models.ToDomainUser(user), err
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

// Implement all methods defined in interfaces.UserRepository...
