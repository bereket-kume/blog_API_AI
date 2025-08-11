package main

import (
	"blog-api/Delivery/routers"
	"blog-api/Infrastructure/database"
	"blog-api/Infrastructure/repositories"
	"blog-api/Infrastructure/services"
	"blog-api/Infrastructure/utils"
	"blog-api/usecases"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	utils.LoadEnv()

	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to MongoDB
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer database.Disconnect()

	// Initialize repositories
	userCollection := database.GetCollection("users")
	blogCollection := database.GetCollection("blogs")
	tokenCollection := database.GetCollection("tokens")

	userRepo := repositories.NewUserMongoRepo(userCollection)
	blogRepo := repositories.NewBlogMongoRepo(blogCollection)
	tokenRepo := repositories.NewTokenMongoRepo(tokenCollection)

	// Initialize recommendation repository
	recommendationRepo := repositories.NewRecommendationMongoRepo(database.GetClient(), database.GetDatabase())

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
	}

	jwtService := services.NewJWTService(jwtSecret, jwtSecret, 15*time.Minute, 7*24*time.Hour)
	passwordService := &services.BcryptHasher{}

	// Initialize email service - Brevo SMTP
	smtpHost := os.Getenv("BREVO_SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp-relay.brevo.com"
	}
	smtpPort := os.Getenv("BREVO_SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}
	smtpUsername := os.Getenv("BREVO_SMTP_USERNAME")
	if smtpUsername == "" {
		smtpUsername = "dummy-key-for-development"
	}
	smtpPassword := os.Getenv("BREVO_SMTP_PASSWORD")
	if smtpPassword == "" {
		smtpPassword = "dummy-key-for-development"
	}
	fromEmail := os.Getenv("FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "noreply@blog-api.com"
	}
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// Log email service configuration
	log.Printf("Email service configuration - Host: %s, Port: %s, Username: %s, From: %s",
		smtpHost, smtpPort, smtpUsername, fromEmail)

	emailService := services.NewEmailService(smtpHost, smtpPort, smtpUsername, smtpPassword, fromEmail, frontendURL)

	// Initialize recommendation service
	recommendationService := services.NewRecommendationService(recommendationRepo, blogRepo)

	// Initialize use cases
	userUC := usecases.NewUserUsecase(userRepo, passwordService, jwtService, tokenRepo, emailService)
	blogUC := usecases.NewBlogUseCase(blogRepo)
	recommendationUC := usecases.NewRecommendationUseCase(recommendationRepo, blogRepo, recommendationService)

	// Create Gin router with proper configuration
	r := gin.New() // Use gin.New() instead of gin.Default() to avoid middleware duplication

	// Add middleware manually
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Configure proxy trust for security
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"}) // Trust only localhost

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Blog API is running",
		})
	})

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Blog API with AI",
			"version": "1.0.0",
		})
	})

	// Initialize recommendation worker
	recommendationWorker := services.NewRecommendationWorker(recommendationUC)
	recommendationWorker.Start()
	defer recommendationWorker.Stop()

	// Setup routes
	routers.SetupRouter(r, userUC, blogUC, recommendationUC, jwtService)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatal("Invalid PORT environment variable")
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
