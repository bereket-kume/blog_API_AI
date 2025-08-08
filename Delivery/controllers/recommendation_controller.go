package controllers

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecommendationController struct {
	recommendationUC interfaces.RecommendationUseCase
}

func NewRecommendationController(recommendationUC interfaces.RecommendationUseCase) *RecommendationController {
	return &RecommendationController{
		recommendationUC: recommendationUC,
	}
}

// TrackUserAction tracks a user action (view, like, comment, etc.)
func (rc *RecommendationController) TrackUserAction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request struct {
		BlogID string `json:"blog_id" binding:"required"`
		Action string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate action
	validActions := []string{models.ActionView, models.ActionLike, models.ActionComment, models.ActionShare, models.ActionBookmark}
	isValidAction := false
	for _, action := range validActions {
		if request.Action == action {
			isValidAction = true
			break
		}
	}

	if !isValidAction {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	err := rc.recommendationUC.TrackUserAction(userID.(string), request.BlogID, request.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track user action"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Action tracked successfully"})
}

// GetUserRecommendations gets personalized recommendations for the user
func (rc *RecommendationController) GetUserRecommendations(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	category := c.DefaultQuery("category", models.CategoryAll)

	response, err := rc.recommendationUC.GetUserRecommendations(userID.(string), limit, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSimilarContent gets content similar to a specific blog
func (rc *RecommendationController) GetSimilarContent(c *gin.Context) {
	blogID := c.Param("id")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 20 {
		limit = 5
	}

	similarBlogs, err := rc.recommendationUC.GetSimilarContent(blogID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get similar content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"blog_id":         blogID,
		"similar_content": similarBlogs,
		"count":           len(similarBlogs),
	})
}

// GetTrendingContent gets currently trending content
func (rc *RecommendationController) GetTrendingContent(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	trendingBlogs, err := rc.recommendationUC.GetTrendingContent(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trending content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trending_content": trendingBlogs,
		"count":            len(trendingBlogs),
	})
}

// GetPopularContent gets popular content
func (rc *RecommendationController) GetPopularContent(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	popularBlogs, err := rc.recommendationUC.GetPopularContent(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get popular content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"popular_content": popularBlogs,
		"count":           len(popularBlogs),
	})
}

// GetNewContent gets recently published content
func (rc *RecommendationController) GetNewContent(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	newBlogs, err := rc.recommendationUC.GetNewContent(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get new content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"new_content": newBlogs,
		"count":       len(newBlogs),
	})
}

// GetUserInterests gets user's interest profile
func (rc *RecommendationController) GetUserInterests(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	interests, err := rc.recommendationUC.GetUserInterests(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user interests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"interests": interests,
		"count":     len(interests),
	})
}

// GetUserBehaviorSummary gets a summary of user's behavior
func (rc *RecommendationController) GetUserBehaviorSummary(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	summary, err := rc.recommendationUC.GetUserBehaviorSummary(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get behavior summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"summary": summary,
	})
}

// GetRecommendationStats gets recommendation statistics for the user
func (rc *RecommendationController) GetRecommendationStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	stats, err := rc.recommendationUC.GetRecommendationStats(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendation stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// MarkRecommendationViewed marks a recommendation as viewed
func (rc *RecommendationController) MarkRecommendationViewed(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	recommendationID := c.Param("id")
	if recommendationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recommendation ID is required"})
		return
	}

	err := rc.recommendationUC.MarkRecommendationViewed(recommendationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark recommendation as viewed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recommendation marked as viewed"})
}

// GetContentDiscovery gets various types of content for discovery
func (rc *RecommendationController) GetContentDiscovery(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 20 {
		limit = 5
	}

	// Get different types of content concurrently
	type contentResult struct {
		Type  string
		Blogs []models.Blog
		Error error
	}

	results := make(chan contentResult, 4)

	// Trending content
	go func() {
		blogs, err := rc.recommendationUC.GetTrendingContent(limit)
		results <- contentResult{"trending", blogs, err}
	}()

	// Popular content
	go func() {
		blogs, err := rc.recommendationUC.GetPopularContent(limit)
		results <- contentResult{"popular", blogs, err}
	}()

	// New content
	go func() {
		blogs, err := rc.recommendationUC.GetNewContent(limit)
		results <- contentResult{"new", blogs, err}
	}()

	// Collect results
	discovery := make(map[string]interface{})
	for i := 0; i < 3; i++ {
		result := <-results
		if result.Error != nil {
			discovery[result.Type] = gin.H{"error": "Failed to fetch " + result.Type + " content"}
		} else {
			discovery[result.Type] = gin.H{
				"content": result.Blogs,
				"count":   len(result.Blogs),
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"discovery": discovery,
	})
}
