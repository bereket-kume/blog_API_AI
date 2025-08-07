package main

import (
	"blog-api/Delivery/routers"
	"blog-api/Infrastructure/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	database.Connect()

	// Set up the Gin router with all dependencies
	router := routers.SetupRouter()

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
