package main

import (
	"blog-api/Delivery/routers"
	"blog-api/Infrastructure/database"
	"blog-api/Infrastructure/repositories"
	"blog-api/Infrastructure/services"
	"blog-api/usecases"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to MongoDB
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer database.Disconnect()

	// Initialize MongoDB collections
	userCollection := database.GetCollection("users")
	blogCollection := database.GetCollection("blogs")
	tokenCollection := database.GetCollection("tokens")
	aiSuggestionCollection := database.GetCollection("ai_suggestions")

	// Initialize repositories
	userRepo := repositories.NewUserMongoRepo(userCollection)
	blogRepo := repositories.NewBlogMongoRepo(blogCollection)
	tokenRepo := repositories.NewTokenMongoRepo(tokenCollection)
	aiSuggestionRepo := repositories.NewAISuggestionMongoRepo(aiSuggestionCollection)

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
	}
	jwtService := services.NewJWTService(jwtSecret, jwtSecret, 15*time.Minute, 7*24*time.Hour)
	passwordService := &services.BcryptHasher{}

	// Initialize use cases
	userUC := usecases.NewUserUsecase(userRepo, passwordService, jwtService, tokenRepo)
	blogUC := usecases.NewBlogUseCase(blogRepo)
	aiSuggestionUC := usecases.NewAISuggestionUseCase(aiSuggestionRepo, blogRepo)

	// Create Gin router
	r := gin.Default()

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

	// Setup routes with use cases
	routers.SetupRouter(r, userUC, blogUC, aiSuggestionUC, jwtService)

	// Run server on port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatal("Invalid PORT environment variable")
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
