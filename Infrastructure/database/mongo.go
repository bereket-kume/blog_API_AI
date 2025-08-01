package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
}

func ConnectDB() (*MongoDB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoUrl := os.Getenv("URL")
	if mongoUrl == "" {
		log.Fatal("mongoUrl not set in .env file")
	}
	clientOptions := options.Client().ApplyURI(mongoUrl)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("mongodb connection error: %v", err)
	}

	fmt.Println("mongodb connected")
	log.Println("mongodb connected")
	return &MongoDB{Client: client}, nil
}

func (db *MongoDB) GetCollection(databaseName, collectionName string) *mongo.Collection {
	return db.Client.Database(databaseName).Collection(collectionName)
}
