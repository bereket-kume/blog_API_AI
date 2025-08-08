package services

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"math"
	"sort"
	"strings"
	"time"
)

type recommendationService struct {
	recommendationRepo interfaces.RecommendationRepository
	blogRepo           interfaces.BlogRepository
}

func NewRecommendationService(
	recommendationRepo interfaces.RecommendationRepository,
	blogRepo interfaces.BlogRepository,
) interfaces.RecommendationService {
	return &recommendationService{
		recommendationRepo: recommendationRepo,
		blogRepo:           blogRepo,
	}
}

// TrackUserAction tracks user interactions with content
func (r *recommendationService) TrackUserAction(userID, blogID, action string) error {
	behavior := models.UserBehavior{
		UserID:    userID,
		BlogID:    blogID,
		Action:    action,
		Weight:    getActionWeight(action),
		CreatedAt: time.Now(),
	}
	return r.recommendationRepo.TrackUserBehavior(behavior)
}

// GetUserBehaviorSummary gets a summary of user's behavior
func (r *recommendationService) GetUserBehaviorSummary(userID string) (map[string]interface{}, error) {
	behaviors, err := r.recommendationRepo.GetUserBehaviors(userID, 100)
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"total_actions": len(behaviors),
		"actions":       make(map[string]int),
		"recent_blogs":  make([]string, 0),
		"top_tags":      make([]string, 0),
	}

	blogIDs := make(map[string]bool)
	tagCounts := make(map[string]int)

	for _, behavior := range behaviors {
		// Count actions
		summary["actions"].(map[string]int)[behavior.Action]++

		// Collect unique blog IDs
		blogIDs[behavior.BlogID] = true

		// Get blog details for tag analysis
		blog, err := r.blogRepo.GetBlogByID(behavior.BlogID)
		if err == nil {
			for _, tag := range blog.Tags {
				tagCounts[tag]++
			}
		}
	}

	// Get recent blogs
	recentBlogs := make([]string, 0)
	for blogID := range blogIDs {
		recentBlogs = append(recentBlogs, blogID)
	}
	summary["recent_blogs"] = recentBlogs

	// Get top tags
	type tagCount struct {
		tag   string
		count int
	}
	var tagCountsList []tagCount
	for tag, count := range tagCounts {
		tagCountsList = append(tagCountsList, tagCount{tag, count})
	}
	sort.Slice(tagCountsList, func(i, j int) bool {
		return tagCountsList[i].count > tagCountsList[j].count
	})

	topTags := make([]string, 0)
	for i := 0; i < 10 && i < len(tagCountsList); i++ {
		topTags = append(topTags, tagCountsList[i].tag)
	}
	summary["top_tags"] = topTags

	return summary, nil
}

// CalculateSimilarity calculates similarity between two blogs
func (r *recommendationService) CalculateSimilarity(blog1, blog2 models.Blog) float64 {
	similarity := 0.0
	factors := 0

	// Tag similarity (Jaccard similarity)
	if len(blog1.Tags) > 0 && len(blog2.Tags) > 0 {
		tagSimilarity := calculateTagSimilarity(blog1.Tags, blog2.Tags)
		similarity += tagSimilarity * 0.4 // 40% weight
		factors++
	}

	// Author similarity
	if blog1.AuthorID == blog2.AuthorID {
		similarity += 0.3 // 30% weight
		factors++
	}

	// Content similarity (simple word overlap)
	contentSimilarity := calculateContentSimilarity(blog1.Content, blog2.Content)
	similarity += contentSimilarity * 0.2 // 20% weight
	factors++

	// Category similarity (if we had categories)
	// For now, we'll use a simple approach based on title similarity
	titleSimilarity := calculateTitleSimilarity(blog1.Title, blog2.Title)
	similarity += titleSimilarity * 0.1 // 10% weight
	factors++

	if factors == 0 {
		return 0.0
	}

	return similarity / float64(factors)
}

// FindSimilarContent finds content similar to a given blog
func (r *recommendationService) FindSimilarContent(blogID string, limit int) ([]models.Blog, error) {
	// Get the source blog
	sourceBlog, err := r.blogRepo.GetBlogByID(blogID)
	if err != nil {
		return nil, err
	}

	// Get all published blogs
	allBlogs, err := r.blogRepo.GetPaginatedBlogs(1, 1000) // Get a large number
	if err != nil {
		return nil, err
	}

	// Calculate similarities
	type blogSimilarity struct {
		blog       models.Blog
		similarity float64
	}

	var similarities []blogSimilarity
	for _, blog := range allBlogs {
		if blog.ID == blogID {
			continue // Skip the source blog
		}

		similarity := r.CalculateSimilarity(sourceBlog, blog)
		if similarity > 0.1 { // Only include if similarity > 10%
			similarities = append(similarities, blogSimilarity{blog, similarity})
		}
	}

	// Sort by similarity
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].similarity > similarities[j].similarity
	})

	// Return top similar blogs
	result := make([]models.Blog, 0, limit)
	for i := 0; i < limit && i < len(similarities); i++ {
		result = append(result, similarities[i].blog)
	}

	return result, nil
}

// UpdateUserInterests updates user's interest profile based on their behavior
func (r *recommendationService) UpdateUserInterests(userID string) error {
	behaviors, err := r.recommendationRepo.GetUserBehaviors(userID, 100)
	if err != nil {
		return err
	}

	// Calculate interest weights
	interestWeights := make(map[string]float64)
	totalWeight := 0.0

	for _, behavior := range behaviors {
		blog, err := r.blogRepo.GetBlogByID(behavior.BlogID)
		if err != nil {
			continue
		}

		// Calculate time decay (older behaviors have less weight)
		timeDecay := calculateTimeDecay(behavior.CreatedAt)
		adjustedWeight := behavior.Weight * timeDecay

		// Add weight to tags
		for _, tag := range blog.Tags {
			interestWeights[tag] += adjustedWeight
		}

		// Add weight to author
		authorKey := "author:" + blog.AuthorID
		interestWeights[authorKey] += adjustedWeight

		totalWeight += adjustedWeight
	}

	// Normalize and create interest records
	for topic, weight := range interestWeights {
		if totalWeight > 0 {
			normalizedWeight := weight / totalWeight
			if normalizedWeight > 0.01 { // Only track interests > 1%
				interest := models.UserInterest{
					UserID:    userID,
					Topic:     topic,
					Weight:    normalizedWeight,
					LastSeen:  time.Now(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				r.recommendationRepo.UpdateUserInterest(interest)
			}
		}
	}

	return nil
}

// GetUserInterestProfile gets user's interest profile
func (r *recommendationService) GetUserInterestProfile(userID string) ([]models.UserInterest, error) {
	return r.recommendationRepo.GetUserInterests(userID)
}

// GenerateUserRecommendations generates personalized recommendations for a user
func (r *recommendationService) GenerateUserRecommendations(userID string, limit int) ([]models.UserRecommendation, error) {
	// Get user interests
	interests, err := r.recommendationRepo.GetUserInterests(userID)
	if err != nil {
		return nil, err
	}

	// Get all published blogs
	allBlogs, err := r.blogRepo.GetPaginatedBlogs(1, 1000)
	if err != nil {
		return nil, err
	}

	// Calculate recommendation scores
	type blogScore struct {
		blog     models.Blog
		score    float64
		reason   string
		category string
	}

	var scoredBlogs []blogScore
	for _, blog := range allBlogs {
		score := 0.0
		reasons := make([]string, 0)

		// Score based on user interests
		for _, interest := range interests {
			if strings.Contains(interest.Topic, "author:") {
				// Author interest
				authorID := strings.TrimPrefix(interest.Topic, "author:")
				if blog.AuthorID == authorID {
					score += interest.Weight * 2.0
					reasons = append(reasons, "Based on your interest in this author")
				}
			} else {
				// Tag interest
				for _, tag := range blog.Tags {
					if tag == interest.Topic {
						score += interest.Weight * 1.5
						reasons = append(reasons, "Based on your interest in "+tag)
						break
					}
				}
			}
		}

		// Add popularity bonus
		popularityScore := float64(blog.ViewCount+blog.Likes*2) / 100.0
		score += popularityScore * 0.3

		// Add recency bonus
		daysSinceCreation := time.Since(blog.CreatedAt).Hours() / 24
		recencyScore := math.Max(0, 1.0-daysSinceCreation/30.0) // Decay over 30 days
		score += recencyScore * 0.2

		if score > 0.1 { // Only include if score > 10%
			reason := "Recommended based on your interests"
			if len(reasons) > 0 {
				reason = reasons[0]
			}

			category := models.CategoryBasedOnLikes
			if popularityScore > 0.5 {
				category = models.CategoryPopular
			}

			scoredBlogs = append(scoredBlogs, blogScore{
				blog:     blog,
				score:    score,
				reason:   reason,
				category: category,
			})
		}
	}

	// Sort by score
	sort.Slice(scoredBlogs, func(i, j int) bool {
		return scoredBlogs[i].score > scoredBlogs[j].score
	})

	// Create recommendation records
	recommendations := make([]models.UserRecommendation, 0, limit)
	for i := 0; i < limit && i < len(scoredBlogs); i++ {
		scored := scoredBlogs[i]
		recommendation := models.UserRecommendation{
			UserID:      userID,
			BlogID:      scored.blog.ID,
			Score:       scored.score,
			Reason:      scored.reason,
			Category:    scored.category,
			GeneratedAt: time.Now(),
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour), // Expire in 7 days
			IsViewed:    false,
		}
		recommendations = append(recommendations, recommendation)
	}

	// Save recommendations
	for _, rec := range recommendations {
		r.recommendationRepo.CreateUserRecommendation(rec)
	}

	return recommendations, nil
}

// GetRecommendations gets recommendations for a user
func (r *recommendationService) GetRecommendations(request models.RecommendationRequest) (models.RecommendationResponse, error) {
	recommendations, err := r.recommendationRepo.GetUserRecommendations(request.UserID, request.Limit, request.Category)
	if err != nil {
		return models.RecommendationResponse{}, err
	}

	blogRecommendations := make([]models.BlogRecommendation, 0, len(recommendations))
	for _, rec := range recommendations {
		blog, err := r.blogRepo.GetBlogByID(rec.BlogID)
		if err != nil {
			continue
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
		UserID:          request.UserID,
		Recommendations: blogRecommendations,
		GeneratedAt:     time.Now(),
		TotalCount:      len(blogRecommendations),
	}, nil
}

// MarkRecommendationViewed marks a recommendation as viewed
func (r *recommendationService) MarkRecommendationViewed(recommendationID string) error {
	return r.recommendationRepo.UpdateRecommendationViewed(recommendationID)
}

// GetTrendingContent gets currently trending content
func (r *recommendationService) GetTrendingContent(limit int) ([]models.Blog, error) {
	return r.recommendationRepo.GetTrendingBlogs(limit)
}

// GetPopularContent gets popular content
func (r *recommendationService) GetPopularContent(limit int) ([]models.Blog, error) {
	// Get all blogs and sort by popularity
	allBlogs, err := r.blogRepo.GetPaginatedBlogs(1, 1000)
	if err != nil {
		return nil, err
	}

	// Sort by popularity score (views + likes*2 + comments*3)
	sort.Slice(allBlogs, func(i, j int) bool {
		scoreI := float64(allBlogs[i].ViewCount + allBlogs[i].Likes*2 + len(allBlogs[i].Comments)*3)
		scoreJ := float64(allBlogs[j].ViewCount + allBlogs[j].Likes*2 + len(allBlogs[j].Comments)*3)
		return scoreI > scoreJ
	})

	if limit > len(allBlogs) {
		limit = len(allBlogs)
	}

	return allBlogs[:limit], nil
}

// GetNewContent gets recently published content
func (r *recommendationService) GetNewContent(limit int) ([]models.Blog, error) {
	// Get all blogs and sort by creation date
	allBlogs, err := r.blogRepo.GetPaginatedBlogs(1, 1000)
	if err != nil {
		return nil, err
	}

	// Sort by creation date (newest first)
	sort.Slice(allBlogs, func(i, j int) bool {
		return allBlogs[i].CreatedAt.After(allBlogs[j].CreatedAt)
	})

	if limit > len(allBlogs) {
		limit = len(allBlogs)
	}

	return allBlogs[:limit], nil
}

// ProcessContentSimilarities processes content similarity calculations
func (r *recommendationService) ProcessContentSimilarities() error {
	// This would be implemented to calculate similarities between all content
	// For now, we'll return nil as this is a background process
	return nil
}

// ProcessUserRecommendations processes user recommendations
func (r *recommendationService) ProcessUserRecommendations() error {
	// This would be implemented to generate recommendations for all users
	// For now, we'll return nil as this is a background process
	return nil
}

// CleanupOldData cleans up old recommendation data
func (r *recommendationService) CleanupOldData() error {
	// Clean up old behaviors (older than 90 days)
	err := r.recommendationRepo.CleanupOldBehaviors(90)
	if err != nil {
		return err
	}

	// Clean up old similarities (older than 30 days)
	err = r.recommendationRepo.CleanupOldSimilarities(30)
	if err != nil {
		return err
	}

	// Delete expired recommendations
	return r.recommendationRepo.DeleteExpiredRecommendations()
}

// GetRecommendationAnalytics gets recommendation statistics for a user
func (r *recommendationService) GetRecommendationAnalytics(userID string) (models.RecommendationStats, error) {
	return r.recommendationRepo.GetRecommendationStats(userID)
}

// GetSystemRecommendationStats gets system-wide recommendation statistics
func (r *recommendationService) GetSystemRecommendationStats() (map[string]interface{}, error) {
	// This would return system-wide statistics
	// For now, return a simple structure
	return map[string]interface{}{
		"total_recommendations_generated": 0,
		"active_users":                    0,
		"average_recommendation_score":    0.0,
		"last_processed_at":               time.Now(),
	}, nil
}

// Helper functions

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

func calculateTagSimilarity(tags1, tags2 []string) float64 {
	if len(tags1) == 0 && len(tags2) == 0 {
		return 1.0
	}
	if len(tags1) == 0 || len(tags2) == 0 {
		return 0.0
	}

	tagSet1 := make(map[string]bool)
	tagSet2 := make(map[string]bool)

	for _, tag := range tags1 {
		tagSet1[tag] = true
	}
	for _, tag := range tags2 {
		tagSet2[tag] = true
	}

	intersection := 0
	for tag := range tagSet1 {
		if tagSet2[tag] {
			intersection++
		}
	}

	union := len(tagSet1) + len(tagSet2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

func calculateContentSimilarity(content1, content2 string) float64 {
	words1 := strings.Fields(strings.ToLower(content1))
	words2 := strings.Fields(strings.ToLower(content2))

	if len(words1) == 0 && len(words2) == 0 {
		return 1.0
	}
	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	wordSet1 := make(map[string]bool)
	wordSet2 := make(map[string]bool)

	for _, word := range words1 {
		if len(word) > 3 { // Only consider words longer than 3 characters
			wordSet1[word] = true
		}
	}
	for _, word := range words2 {
		if len(word) > 3 {
			wordSet2[word] = true
		}
	}

	intersection := 0
	for word := range wordSet1 {
		if wordSet2[word] {
			intersection++
		}
	}

	union := len(wordSet1) + len(wordSet2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

func calculateTitleSimilarity(title1, title2 string) float64 {
	return calculateContentSimilarity(title1, title2)
}

func calculateTimeDecay(createdAt time.Time) float64 {
	daysSince := time.Since(createdAt).Hours() / 24
	// Exponential decay with half-life of 30 days
	return math.Exp(-daysSince / 30.0 * math.Ln2)
}
