package main

import (
	"blog-api/Delivery/routers"
	Database "blog-api/Infrastructure/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Database.Connect()

	r := routers.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
