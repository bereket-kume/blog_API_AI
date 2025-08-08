package interfaces

import "blog-api/Domain/models"

// RecommendationService defines the interface for recommendation business logic
type RecommendationService interface {
	// User Behavior Tracking
	TrackUserAction(userID, blogID, action string) error
	GetUserBehaviorSummary(userID string) (map[string]interface{}, error)

	// Content Similarity
	CalculateSimilarity(blog1, blog2 models.Blog) float64
	FindSimilarContent(blogID string, limit int) ([]models.Blog, error)

	// User Interests
	UpdateUserInterests(userID string) error
	GetUserInterestProfile(userID string) ([]models.UserInterest, error)

	// Recommendation Generation
	GenerateUserRecommendations(userID string, limit int) ([]models.UserRecommendation, error)
	GetRecommendations(request models.RecommendationRequest) (models.RecommendationResponse, error)
	MarkRecommendationViewed(recommendationID string) error

	// Content Analysis
	GetTrendingContent(limit int) ([]models.Blog, error)
	GetPopularContent(limit int) ([]models.Blog, error)
	GetNewContent(limit int) ([]models.Blog, error)

	// Background Processing
	ProcessContentSimilarities() error
	ProcessUserRecommendations() error
	CleanupOldData() error

	// Analytics
	GetRecommendationAnalytics(userID string) (models.RecommendationStats, error)
	GetSystemRecommendationStats() (map[string]interface{}, error)
}
