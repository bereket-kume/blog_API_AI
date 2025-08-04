package Database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var client *mongo.Client

func Connect() {
	uri := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Correct: pass context first, then options
	clientOptions := options.Client().ApplyURI(uri)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	// Ping to confirm connection
	err = client.Ping(ctx, nil)
	if err != nil {
		panic("MongoDB ping error: " + err.Error())
	}

	dbName := os.Getenv("MONGO_DB")
	db = client.Database(dbName)
	fmt.Println("✅ MongoDB connected to database:", dbName)
}

func GetUserCollection() *mongo.Collection {
	return db.Collection("users")
}
