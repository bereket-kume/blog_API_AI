# 🎯 **Recommendation System for Blog API**

A comprehensive recommendation system that works in the background to provide personalized content recommendations based on user behavior and content similarity.

## 🚀 **Features**

### **Core Functionality**
- **User Behavior Tracking** - Track views, likes, comments, shares, and bookmarks
- **Content Similarity** - Find related content based on tags, authors, and content analysis
- **Personalized Recommendations** - Generate recommendations based on user interests
- **Background Processing** - Async recommendation generation and similarity calculations
- **Content Discovery** - Trending, popular, and new content discovery

### **Recommendation Types**
- **Based on Likes** - Content similar to what the user has liked
- **Similar Content** - Content similar to what the user is currently viewing
- **Trending Content** - Currently popular content
- **Popular Content** - High-engagement content
- **New Content** - Recently published content

## 🏗️ **Architecture**

### **Clean Architecture Layers**

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Controllers   │  │    Routers      │  │  Middleware  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                   Use Case Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │ Recommendation  │  │   Blog Use      │  │  User Use    │ │
│  │    Use Case     │  │     Case        │  │   Case       │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                 Service Layer                               │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │ Recommendation  │  │   Background    │  │  Analytics   │ │
│  │    Service      │  │    Worker       │  │   Service    │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│               Infrastructure Layer                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │ Recommendation  │  │   Blog Repo     │  │  User Repo   │ │
│  │   Repository    │  │                 │  │              │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │    Models       │  │   Interfaces    │  │  Constants   │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 📊 **Data Models**

### **User Behavior**
```go
type UserBehavior struct {
    ID        string    `json:"id" bson:"_id,omitempty"`
    UserID    string    `json:"user_id" bson:"user_id"`
    BlogID    string    `json:"blog_id" bson:"blog_id"`
    Action    string    `json:"action" bson:"action"` // view, like, comment, share, bookmark
    Weight    float64   `json:"weight" bson:"weight"` // importance of this action
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
```

### **User Recommendation**
```go
type UserRecommendation struct {
    ID          string    `json:"id" bson:"_id,omitempty"`
    UserID      string    `json:"user_id" bson:"user_id"`
    BlogID      string    `json:"blog_id" bson:"blog_id"`
    Score       float64   `json:"score" bson:"score"`
    Reason      string    `json:"reason" bson:"reason"`
    Category    string    `json:"category" bson:"category"`
    GeneratedAt time.Time `json:"generated_at" bson:"generated_at"`
    ExpiresAt   time.Time `json:"expires_at" bson:"expires_at"`
    IsViewed    bool      `json:"is_viewed" bson:"is_viewed"`
    ViewedAt    *time.Time `json:"viewed_at" bson:"viewed_at,omitempty"`
}
```

### **User Interest**
```go
type UserInterest struct {
    ID        string    `json:"id" bson:"_id,omitempty"`
    UserID    string    `json:"user_id" bson:"user_id"`
    Topic     string    `json:"topic" bson:"topic"`         // tag, category, author
    Weight    float64   `json:"weight" bson:"weight"`       // interest strength (0.0 to 1.0)
    LastSeen  time.Time `json:"last_seen" bson:"last_seen"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
```

## 🔧 **API Endpoints**

### **Public Endpoints**

#### **Content Discovery**
```http
GET /recommendations/trending?limit=10
GET /recommendations/popular?limit=10
GET /recommendations/new?limit=10
GET /recommendations/discovery?limit=5
```

#### **Similar Content**
```http
GET /blogs/{id}/similar?limit=5
```

### **Authenticated Endpoints**

#### **User Actions**
```http
POST /api/recommendations/track
Content-Type: application/json

{
    "blog_id": "blog123",
    "action": "view"  // view, like, comment, share, bookmark
}
```

#### **Personalized Recommendations**
```http
GET /api/recommendations/personal?limit=10&category=all
GET /api/recommendations/interests
GET /api/recommendations/behavior
GET /api/recommendations/stats
```

#### **Recommendation Tracking**
```http
PUT /api/recommendations/{id}/view
```

## 🎯 **How It Works**

### **1. User Behavior Tracking**
When users interact with content, their actions are tracked:
- **View** (Weight: 1.0) - User viewed a blog
- **Like** (Weight: 5.0) - User liked a blog
- **Comment** (Weight: 3.0) - User commented on a blog
- **Share** (Weight: 4.0) - User shared a blog
- **Bookmark** (Weight: 2.0) - User bookmarked a blog

### **2. Interest Calculation**
User interests are calculated based on:
- **Tags** - What topics the user engages with
- **Authors** - Which authors the user follows
- **Time Decay** - Recent actions have more weight
- **Action Weight** - Different actions have different importance

### **3. Content Similarity**
Content similarity is calculated using:
- **Tag Overlap** (40% weight) - Jaccard similarity of tags
- **Author Similarity** (30% weight) - Same author
- **Content Similarity** (20% weight) - Word overlap analysis
- **Title Similarity** (10% weight) - Title word overlap

### **4. Recommendation Generation**
Recommendations are generated using:
- **Interest-based scoring** - Based on user's interest profile
- **Popularity bonus** - High-engagement content gets bonus
- **Recency bonus** - New content gets bonus
- **Category classification** - Based on recommendation reason

### **5. Background Processing**
The system runs background tasks:
- **Content Similarity Updates** - Calculate similarities between content
- **User Recommendation Generation** - Generate recommendations for users
- **Data Cleanup** - Remove old behavior and similarity data

## 🔄 **Background Worker**

The recommendation system includes a background worker that:
- Runs every hour by default
- Processes content similarities
- Generates user recommendations
- Cleans up old data
- Can be started/stopped programmatically

```go
// Start the worker
recommendationWorker := services.NewRecommendationWorker(recommendationUC)
recommendationWorker.Start()
defer recommendationWorker.Stop()
```

## 📈 **Analytics & Insights**

### **User Behavior Summary**
```json
{
    "user_id": "user123",
    "summary": {
        "total_actions": 45,
        "actions": {
            "view": 30,
            "like": 10,
            "comment": 3,
            "share": 2
        },
        "recent_blogs": ["blog1", "blog2", "blog3"],
        "top_tags": ["go", "programming", "web-development"]
    }
}
```

### **Recommendation Statistics**
```json
{
    "user_id": "user123",
    "total_recommendations": 150,
    "viewed_recommendations": 45,
    "clicked_recommendations": 12,
    "average_score": 0.75,
    "last_generated_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
}
```

## 🚀 **Getting Started**

### **1. Prerequisites**
- Go 1.19+
- MongoDB (local or cloud)
- Existing blog API setup

### **2. Integration**
The recommendation system is already integrated into the main application. Just start the server:

```bash
go run ./Delivery/main.go
```

### **3. Usage Examples**

#### **Track User Action**
```bash
curl -X POST http://localhost:8080/api/recommendations/track \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"blog_id": "blog123", "action": "like"}'
```

#### **Get Personal Recommendations**
```bash
curl -X GET "http://localhost:8080/api/recommendations/personal?limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### **Get Similar Content**
```bash
curl -X GET "http://localhost:8080/blogs/blog123/similar?limit=5"
```

#### **Get Trending Content**
```bash
curl -X GET "http://localhost:8080/recommendations/trending?limit=10"
```

## 🔧 **Configuration**

### **Environment Variables**
```bash
# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
DB_NAME=blog_db

# Recommendation System
RECOMMENDATION_WORKER_INTERVAL=1h
RECOMMENDATION_EXPIRY_DAYS=7
BEHAVIOR_CLEANUP_DAYS=90
SIMILARITY_CLEANUP_DAYS=30
```

### **Action Weights**
```go
const (
    WeightView    = 1.0
    WeightLike    = 5.0
    WeightComment = 3.0
    WeightShare   = 4.0
    WeightBookmark = 2.0
)
```

## 📊 **Performance Considerations**

### **Optimizations**
- **Background Processing** - Heavy calculations run asynchronously
- **Caching** - Recommendations are cached and expire after 7 days
- **Batch Processing** - Similarity calculations are batched
- **Data Cleanup** - Old data is automatically cleaned up

### **Scaling**
- **Horizontal Scaling** - Multiple worker instances can run
- **Database Indexing** - Proper indexes on user_id, blog_id, created_at
- **Connection Pooling** - MongoDB connection pooling
- **Rate Limiting** - API rate limiting for tracking endpoints

## 🧪 **Testing**

### **Unit Tests**
```bash
# Run recommendation unit tests
go test -v ./usecases/ -run TestRecommendation
```

### **Integration Tests**
```bash
# Run recommendation integration tests
go test -v ./Infrastructure/repositories/ -run TestRecommendation
```

### **API Tests**
```bash
# Test recommendation endpoints
curl -X GET "http://localhost:8080/recommendations/trending"
```

## 🔮 **Future Enhancements**

### **Planned Features**
- **Machine Learning Integration** - ML-based recommendation algorithms
- **Real-time Processing** - Real-time recommendation updates
- **A/B Testing** - Test different recommendation strategies
- **Advanced Analytics** - Detailed user behavior analytics
- **Content Clustering** - Advanced content similarity algorithms

### **Advanced Algorithms**
- **Collaborative Filtering** - User-based and item-based filtering
- **Content-Based Filtering** - Advanced content analysis
- **Hybrid Approaches** - Combine multiple recommendation strategies
- **Deep Learning** - Neural network-based recommendations

## 📝 **API Documentation**

For complete API documentation, see the individual endpoint documentation in the codebase or run the server and visit the health endpoint for basic information.

## 🤝 **Contributing**

1. Follow Clean Architecture principles
2. Add tests for new features
3. Update documentation
4. Follow Go coding standards
5. Test with local MongoDB

## 📄 **License**

This recommendation system is part of the Blog API project and follows the same licensing terms. 