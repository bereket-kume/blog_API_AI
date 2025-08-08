package models

import (
	"time"
)

// UserBehavior represents user interactions with content
type UserBehavior struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	BlogID    string    `json:"blog_id" bson:"blog_id"`
	Action    string    `json:"action" bson:"action"` // like, view, comment, share
	Weight    float64   `json:"weight" bson:"weight"` // importance of this action
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// ContentSimilarity represents similarity between content items
type ContentSimilarity struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	BlogID1     string    `json:"blog_id_1" bson:"blog_id_1"`
	BlogID2     string    `json:"blog_id_2" bson:"blog_id_2"`
	Similarity  float64   `json:"similarity" bson:"similarity"` // 0.0 to 1.0
	Factors     []string  `json:"factors" bson:"factors"`       // tags, author, category, etc.
	LastUpdated time.Time `json:"last_updated" bson:"last_updated"`
}

// UserRecommendation represents personalized recommendations for a user
type UserRecommendation struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	UserID      string     `json:"user_id" bson:"user_id"`
	BlogID      string     `json:"blog_id" bson:"blog_id"`
	Score       float64    `json:"score" bson:"score"`       // recommendation score
	Reason      string     `json:"reason" bson:"reason"`     // why this was recommended
	Category    string     `json:"category" bson:"category"` // based_on_likes, similar_content, trending
	GeneratedAt time.Time  `json:"generated_at" bson:"generated_at"`
	ExpiresAt   time.Time  `json:"expires_at" bson:"expires_at"`
	IsViewed    bool       `json:"is_viewed" bson:"is_viewed"`
	ViewedAt    *time.Time `json:"viewed_at" bson:"viewed_at,omitempty"`
}

// RecommendationRequest represents a request for recommendations
type RecommendationRequest struct {
	UserID   string   `json:"user_id"`
	Limit    int      `json:"limit"`
	Category string   `json:"category,omitempty"` // all, based_on_likes, similar_content, trending
	Exclude  []string `json:"exclude,omitempty"`  // blog IDs to exclude
}

// RecommendationResponse represents the response with recommended blogs
type RecommendationResponse struct {
	UserID          string               `json:"user_id"`
	Recommendations []BlogRecommendation `json:"recommendations"`
	GeneratedAt     time.Time            `json:"generated_at"`
	TotalCount      int                  `json:"total_count"`
}

// BlogRecommendation represents a recommended blog with metadata
type BlogRecommendation struct {
	Blog       Blog    `json:"blog"`
	Score      float64 `json:"score"`
	Reason     string  `json:"reason"`
	Category   string  `json:"category"`
	Similarity float64 `json:"similarity,omitempty"` // similarity to user's interests
}

// UserInterest represents user's interest in specific topics
type UserInterest struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Topic     string    `json:"topic" bson:"topic"`   // tag, category, author
	Weight    float64   `json:"weight" bson:"weight"` // interest strength (0.0 to 1.0)
	LastSeen  time.Time `json:"last_seen" bson:"last_seen"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// RecommendationStats represents statistics about recommendations
type RecommendationStats struct {
	UserID                 string    `json:"user_id" bson:"user_id"`
	TotalRecommendations   int       `json:"total_recommendations" bson:"total_recommendations"`
	ViewedRecommendations  int       `json:"viewed_recommendations" bson:"viewed_recommendations"`
	ClickedRecommendations int       `json:"clicked_recommendations" bson:"clicked_recommendations"`
	AverageScore           float64   `json:"average_score" bson:"average_score"`
	LastGeneratedAt        time.Time `json:"last_generated_at" bson:"last_generated_at"`
	UpdatedAt              time.Time `json:"updated_at" bson:"updated_at"`
}

// Action weights for different user behaviors
const (
	ActionView     = "view"
	ActionLike     = "like"
	ActionComment  = "comment"
	ActionShare    = "share"
	ActionBookmark = "bookmark"

	WeightView     = 1.0
	WeightLike     = 5.0
	WeightComment  = 3.0
	WeightShare    = 4.0
	WeightBookmark = 2.0
)

// Recommendation categories
const (
	CategoryBasedOnLikes   = "based_on_likes"
	CategorySimilarContent = "similar_content"
	CategoryTrending       = "trending"
	CategoryPopular        = "popular"
	CategoryNew            = "new"
	CategoryAll            = "all"
)
