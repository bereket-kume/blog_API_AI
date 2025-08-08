package interfaces

import "blog-api/Domain/models"

// RecommendationRepository defines the interface for recommendation data operations
type RecommendationRepository interface {
	// User Behavior Tracking
	TrackUserBehavior(behavior models.UserBehavior) error
	GetUserBehaviors(userID string, limit int) ([]models.UserBehavior, error)
	GetUserBehaviorStats(userID string) (map[string]int, error)

	// Content Similarity
	CalculateContentSimilarity(blogID1, blogID2 string) (float64, error)
	GetSimilarContent(blogID string, limit int) ([]models.ContentSimilarity, error)
	UpdateContentSimilarity(similarity models.ContentSimilarity) error
	GetContentSimilarities(blogID string) ([]models.ContentSimilarity, error)

	// User Recommendations
	CreateUserRecommendation(recommendation models.UserRecommendation) error
	GetUserRecommendations(userID string, limit int, category string) ([]models.UserRecommendation, error)
	UpdateRecommendationViewed(recommendationID string) error
	DeleteExpiredRecommendations() error
	GetRecommendationStats(userID string) (models.RecommendationStats, error)

	// User Interests
	UpdateUserInterest(interest models.UserInterest) error
	GetUserInterests(userID string) ([]models.UserInterest, error)
	GetTopUserInterests(userID string, limit int) ([]models.UserInterest, error)

	// Content Analysis
	GetPopularTags(limit int) ([]string, error)
	GetTrendingBlogs(limit int) ([]models.Blog, error)
	GetPopularAuthors(limit int) ([]string, error)

	// Background Processing
	GetBlogsForSimilarityCalculation(limit int) ([]models.Blog, error)
	GetUsersForRecommendationGeneration(limit int) ([]string, error)

	// Utility Methods
	CleanupOldBehaviors(daysOld int) error
	CleanupOldSimilarities(daysOld int) error
}
