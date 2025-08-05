package repositories

import (
	"blog-api/Domain/models"
	"blog-api/Infrastructure/db_models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type tokenMongoRepo struct {
	collection *mongo.Collection
}

func NewTokenMongoRepo(col *mongo.Collection) *tokenMongoRepo {
	return &tokenMongoRepo{collection: col}
}

func (tr *tokenMongoRepo) CreateToken(token models.Token) error {
	db_token := db_models.FromDomainToken(token)
	_, err := tr.collection.InsertOne(context.TODO(), db_token)
	return err
}
func (tr *tokenMongoRepo) DeleteToken(tokenID string) error {
	_, err := tr.collection.DeleteOne(context.TODO(), bson.M{"id": tokenID})
	return err
}
func (tr *tokenMongoRepo) Update(token models.Token) error {
	db_token := db_models.FromDomainToken(token)
	filter := bson.M{"id": db_token.ID}
	update := bson.M{"$set": db_token}
	_, err := tr.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (tr *tokenMongoRepo) GetToken(tokenID string) (models.Token, error) {
	var token db_models.Token
	filter := bson.M{"id": tokenID}
	err := tr.collection.FindOne(context.TODO(), filter).Decode(&token)
	return db_models.ToDomainToken(token), err

}
