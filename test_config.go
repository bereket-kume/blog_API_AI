package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestConfig holds test configuration
type TestConfig struct {
	MongoURI       string
	DatabaseName   string
	CollectionName string
	Timeout        time.Duration
}

// GetTestConfig returns test configuration with defaults
func GetTestConfig() *TestConfig {
	return &TestConfig{
		MongoURI:       getEnvOrDefault("TEST_MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName:   getEnvOrDefault("TEST_DB_NAME", "blog_test_db"),
		CollectionName: getEnvOrDefault("TEST_COLLECTION_NAME", "blogs"),
		Timeout:        10 * time.Second,
	}
}

// ConnectTestDB connects to test MongoDB instance
func ConnectTestDB(config *TestConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to test MongoDB: %s", config.MongoURI)
	return client, nil
}

// DisconnectTestDB disconnects from test MongoDB instance
func DisconnectTestDB(client *mongo.Client, config *TestConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}

	log.Printf("Disconnected from test MongoDB")
	return nil
}

// CleanupTestDB cleans up test database
func CleanupTestDB(client *mongo.Client, config *TestConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	database := client.Database(config.DatabaseName)
	err := database.Drop(ctx)
	if err != nil {
		return err
	}

	log.Printf("Cleaned up test database: %s", config.DatabaseName)
	return nil
}

// GetTestCollection returns test collection
func GetTestCollection(client *mongo.Client, config *TestConfig) *mongo.Collection {
	return client.Database(config.DatabaseName).Collection(config.CollectionName)
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsMongoRunning checks if MongoDB is running on the configured URI
func IsMongoRunning(config *TestConfig) bool {
	client, err := ConnectTestDB(config)
	if err != nil {
		return false
	}
	defer DisconnectTestDB(client, config)
	return true
}
