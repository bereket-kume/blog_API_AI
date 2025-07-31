package main

import (
	"log"
	"os"

	"blog_API_AI/Infrastructure/database"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

}
