package usecases

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"time"
)

type recommendationUseCase struct {
	recommendationRepo interfaces.RecommendationRepository
	blogRepo           interfaces.BlogRepository
	recommendationSvc  interfaces.RecommendationService
}

func NewRecommendationUseCase(
	recommendationRepo interfaces.RecommendationRepository,
	blogRepo interfaces.BlogRepository,
	recommendationSvc interfaces.RecommendationService,
) interfaces.RecommendationUseCase {
	return &recommendationUseCase{
		recommendationRepo: recommendationRepo,
		blogRepo:           blogRepo,
		recommendationSvc:  recommendationSvc,
	}
}

// TrackUserAction tracks user interactions with content
func (r *recommendationUseCase) TrackUserAction(userID, blogID, action string) error {
	// Create behavior record
	behavior := models.UserBehavior{
		UserID:    userID,
		BlogID:    blogID,
		Action:    action,
		Weight:    getActionWeight(action),
		CreatedAt: time.Now(),
	}

	// Track the behavior
	err := r.recommendationRepo.TrackUserBehavior(behavior)
	if err != nil {
		return err
	}

	// Update user interests in background (async)
	go func() {
		r.recommendationSvc.UpdateUserInterests(userID)
	}()

	return nil
}

// GetUserRecommendations retrieves personalized recommendations for a user
func (r *recommendationUseCase) GetUserRecommendations(userID string, limit int, category string) (models.RecommendationResponse, error) {
	// Get stored recommendations
	recommendations, err := r.recommendationRepo.GetUserRecommendations(userID, limit, category)
	if err != nil {
		return models.RecommendationResponse{}, err
	}

	// If no stored recommendations or they're old, generate new ones
	if len(recommendations) == 0 || time.Since(recommendations[0].GeneratedAt) > 24*time.Hour {
		newRecommendations, err := r.recommendationSvc.GenerateUserRecommendations(userID, limit)
		if err != nil {
			return models.RecommendationResponse{}, err
		}
		recommendations = newRecommendations
	}

	// Convert to response format
	blogRecommendations := make([]models.BlogRecommendation, 0, len(recommendations))
	for _, rec := range recommendations {
		blog, err := r.blogRepo.GetBlogByID(rec.BlogID)
		if err != nil {
			continue // Skip if blog not found
		}

		blogRec := models.BlogRecommendation{
			Blog:     blog,
			Score:    rec.Score,
			Reason:   rec.Reason,
			Category: rec.Category,
		}
		blogRecommendations = append(blogRecommendations, blogRec)
	}

	return models.RecommendationResponse{
		UserID:          userID,
		Recommendations: blogRecommendations,
		GeneratedAt:     time.Now(),
		TotalCount:      len(blogRecommendations),
	}, nil
}

// MarkRecommendationViewed marks a recommendation as viewed
func (r *recommendationUseCase) MarkRecommendationViewed(recommendationID string) error {
	return r.recommendationRepo.UpdateRecommendationViewed(recommendationID)
}

// GetSimilarContent finds content similar to a given blog
func (r *recommendationUseCase) GetSimilarContent(blogID string, limit int) ([]models.Blog, error) {
	return r.recommendationSvc.FindSimilarContent(blogID, limit)
}

// GetTrendingContent gets currently trending content
func (r *recommendationUseCase) GetTrendingContent(limit int) ([]models.Blog, error) {
	return r.recommendationSvc.GetTrendingContent(limit)
}

// GetPopularContent gets popular content
func (r *recommendationUseCase) GetPopularContent(limit int) ([]models.Blog, error) {
	return r.recommendationSvc.GetPopularContent(limit)
}

// GetNewContent gets recently published content
func (r *recommendationUseCase) GetNewContent(limit int) ([]models.Blog, error) {
	return r.recommendationSvc.GetNewContent(limit)
}

// GetUserInterests gets user's interest profile
func (r *recommendationUseCase) GetUserInterests(userID string) ([]models.UserInterest, error) {
	return r.recommendationSvc.GetUserInterestProfile(userID)
}

// GetUserBehaviorSummary gets a summary of user's behavior
func (r *recommendationUseCase) GetUserBehaviorSummary(userID string) (map[string]interface{}, error) {
	return r.recommendationSvc.GetUserBehaviorSummary(userID)
}

// GetRecommendationStats gets recommendation statistics for a user
func (r *recommendationUseCase) GetRecommendationStats(userID string) (models.RecommendationStats, error) {
	return r.recommendationRepo.GetRecommendationStats(userID)
}

// ProcessRecommendations processes recommendations in the background
func (r *recommendationUseCase) ProcessRecommendations() error {
	return r.recommendationSvc.ProcessUserRecommendations()
}

// UpdateContentSimilarities updates content similarity calculations
func (r *recommendationUseCase) UpdateContentSimilarities() error {
	return r.recommendationSvc.ProcessContentSimilarities()
}

// CleanupOldData cleans up old recommendation data
func (r *recommendationUseCase) CleanupOldData() error {
	return r.recommendationSvc.CleanupOldData()
}

// Helper function to get action weight
func getActionWeight(action string) float64 {
	switch action {
	case models.ActionView:
		return models.WeightView
	case models.ActionLike:
		return models.WeightLike
	case models.ActionComment:
		return models.WeightComment
	case models.ActionShare:
		return models.WeightShare
	case models.ActionBookmark:
		return models.WeightBookmark
	default:
		return 1.0
	}
}
