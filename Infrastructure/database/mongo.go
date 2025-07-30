package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
)

// Connect establishes connection to MongoDB Atlas
func Connect() error {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetServerSelectionTimeout(30 * time.Second)
	clientOptions.SetSocketTimeout(30 * time.Second)

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Test the connection
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		log.Println("Please check:")
		log.Println("1. Your IP address is whitelisted in MongoDB Atlas")
		log.Println("2. Database user has correct permissions")
		log.Println("3. Connection string format is correct")
		return err
	}

	// Set the database
	DB = Client.Database("blog_db")
	log.Println("Successfully connected to MongoDB Atlas!")
	return nil
}

// Disconnect closes the MongoDB connection
func Disconnect() error {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return Client.Disconnect(ctx)
	}
	return nil
}

// GetCollection returns a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}

// GetClient returns the MongoDB client
func GetClient() *mongo.Client {
	return Client
}

// GetDatabase returns the MongoDB database
func GetDatabase() *mongo.Database {
	return DB
}
