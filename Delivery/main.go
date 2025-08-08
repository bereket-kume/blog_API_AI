package main

import (
	"blog-api/Delivery/routers"
	"blog-api/Infrastructure/database"
	"blog-api/Infrastructure/repositories"
	"blog-api/Infrastructure/services"
	"blog-api/usecases"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load env vars (optional if you already set MONGODB_URI before running)
	godotenv.Load()

	// Connect to MongoDB
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer database.Disconnect()

	// === Config ===
	accessSecret := os.Getenv("ACCESS_SECRET")
	if accessSecret == "" {
		accessSecret = "supersecretkey"
	}
	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = "anothersecretkey"
	}
	accessTTL := time.Minute * 15
	refreshTTL := time.Hour * 24 * 7

	// === Infrastructure Layer ===
	userRepo := repositories.NewUserMongoRepo(database.GetCollection("users"))
	tokenRepo := repositories.NewTokenMongoRepo(database.GetCollection("tokens"))
	emailService := services.NewEmailService(
		os.Getenv("RESEND_API_KEY"),
		os.Getenv("FROM_EMAIL"),
		os.Getenv("FRONTEND_URL"),
	)

	hasher := services.BcryptHasher{}
	jwtService := services.NewJWTService(accessSecret, refreshSecret, accessTTL, refreshTTL)

	// === Usecases ===
	userUC := usecases.NewUserUsecase(userRepo, hasher, jwtService, tokenRepo, emailService)

	// === Setup Router ===
	r := gin.Default()
	routers.SetupRouter(r, userUC, jwtService)

	// === Run Server ===
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

// package main

// import (
// 	"fmt"

// 	"github.com/resend/resend-go/v2"
// )

// func main() {
// 	apiKey := "re_dBw2b75y_3J2sVAFqtJTgq4HaXx3LTsZ3"

// 	client := resend.NewClient(apiKey)

// 	params := &resend.SendEmailRequest{
// 		From:    "onboarding@resend.dev",
// 		To:      []string{"awelabubekar625@gmail.com"},
// 		Subject: "Hello World",
// 		Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
// 	}

// 	sent, err := client.Emails.Send(params)
// 	fmt.Print(err, sent)
// }
