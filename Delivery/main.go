package main

import (
	"blog-api/Delivery/routers"
	"blog-api/Infrastructure/database"
	"blog-api/Infrastructure/repositories"
	"blog-api/Infrastructure/services"
	"blog-api/usecases"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	// Initialize recommendation repository
	recommendationRepo := repositories.NewRecommendationMongoRepo(database.GetClient(), database.GetDatabase())

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
	}
	jwtService := services.NewJWTService(jwtSecret, jwtSecret, 15*time.Minute, 7*24*time.Hour)
	passwordService := &services.BcryptHasher{}

	// Initialize email service
	emailAPIKey := os.Getenv("RESEND_API_KEY")
	if emailAPIKey == "" {
		emailAPIKey = "dummy-key-for-development"
	}
	fromEmail := os.Getenv("FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "noreply@blog-api.com"
	}
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	emailService := services.NewEmailService(emailAPIKey, fromEmail, frontendURL)

	// Initialize recommendation service
	recommendationService := services.NewRecommendationService(recommendationRepo, blogRepo)

	// Initialize use cases
	userUC := usecases.NewUserUsecase(userRepo, passwordService, jwtService, tokenRepo, emailService)
	blogUC := usecases.NewBlogUseCase(blogRepo)
	recommendationUC := usecases.NewRecommendationUseCase(recommendationRepo, blogRepo, recommendationService)
	aiSuggestionUC := usecases.NewAISuggestionUseCase(aiSuggestionRepo, blogRepo)

	// Create Gin router
	r := gin.Default()

	// Add timeout middleware
	r.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		// Get allowed origins from environment or use default
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "*"
		}

		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

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
	if err := recommendationWorker.Start(); err != nil {
		log.Printf("Warning: Failed to start recommendation worker: %v", err)
	} else {
		defer recommendationWorker.Stop()
	}

	// Setup routes with all use cases
	routers.SetupRouter(r, userUC, blogUC, recommendationUC, aiSuggestionUC, jwtService)

	// Run server on port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatal("Invalid PORT environment variable")
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
