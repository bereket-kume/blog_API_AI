package interfaces

import "blog-api/Domain/models"

// RecommendationUseCase defines the interface for recommendation application logic
type RecommendationUseCase interface {
	// User Actions
	TrackUserAction(userID, blogID, action string) error
	GetUserRecommendations(userID string, limit int, category string) (models.RecommendationResponse, error)
	MarkRecommendationViewed(recommendationID string) error

	// Content Discovery
	GetSimilarContent(blogID string, limit int) ([]models.Blog, error)
	GetTrendingContent(limit int) ([]models.Blog, error)
	GetPopularContent(limit int) ([]models.Blog, error)
	GetNewContent(limit int) ([]models.Blog, error)

	// User Insights
	GetUserInterests(userID string) ([]models.UserInterest, error)
	GetUserBehaviorSummary(userID string) (map[string]interface{}, error)

	// Analytics
	GetRecommendationStats(userID string) (models.RecommendationStats, error)

	// Background Processing (called by background workers)
	ProcessRecommendations() error
	UpdateContentSimilarities() error
	CleanupOldData() error
}
