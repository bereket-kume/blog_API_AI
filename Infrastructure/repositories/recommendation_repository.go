package repositories

import (
	"blog-api/Domain/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type recommendationMongoRepo struct {
	client                    *mongo.Client
	database                  *mongo.Database
	behaviorsCollection       *mongo.Collection
	similaritiesCollection    *mongo.Collection
	recommendationsCollection *mongo.Collection
	interestsCollection       *mongo.Collection
	statsCollection           *mongo.Collection
	blogsCollection           *mongo.Collection
}

func NewRecommendationMongoRepo(client *mongo.Client, database *mongo.Database) *recommendationMongoRepo {
	return &recommendationMongoRepo{
		client:                    client,
		database:                  database,
		behaviorsCollection:       database.Collection("user_behaviors"),
		similaritiesCollection:    database.Collection("content_similarities"),
		recommendationsCollection: database.Collection("user_recommendations"),
		interestsCollection:       database.Collection("user_interests"),
		statsCollection:           database.Collection("recommendation_stats"),
		blogsCollection:           database.Collection("blogs"),
	}
}

// User Behavior Tracking

func (r *recommendationMongoRepo) TrackUserBehavior(behavior models.UserBehavior) error {
	behavior.CreatedAt = time.Now()
	_, err := r.behaviorsCollection.InsertOne(context.TODO(), behavior)
	return err
}

func (r *recommendationMongoRepo) GetUserBehaviors(userID string, limit int) ([]models.UserBehavior, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))

	cursor, err := r.behaviorsCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var behaviors []models.UserBehavior
	if err = cursor.All(context.TODO(), &behaviors); err != nil {
		return nil, err
	}

	return behaviors, nil
}

func (r *recommendationMongoRepo) GetUserBehaviorStats(userID string) (map[string]int, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{
			"_id":   "$action",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := r.behaviorsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	stats := make(map[string]int)
	for cursor.Next(context.TODO()) {
		var result struct {
			Action string `bson:"_id"`
			Count  int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		stats[result.Action] = result.Count
	}

	return stats, nil
}

// Content Similarity

func (r *recommendationMongoRepo) CalculateContentSimilarity(blogID1, blogID2 string) (float64, error) {
	// This would implement a more sophisticated similarity calculation
	// For now, return a simple similarity based on tag overlap
	blog1, err := r.getBlogByID(blogID1)
	if err != nil {
		return 0.0, err
	}

	blog2, err := r.getBlogByID(blogID2)
	if err != nil {
		return 0.0, err
	}

	// Simple tag similarity
	similarity := calculateTagOverlap(blog1.Tags, blog2.Tags)
	return similarity, nil
}

func (r *recommendationMongoRepo) GetSimilarContent(blogID string, limit int) ([]models.ContentSimilarity, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"blog_id_1": blogID},
			{"blog_id_2": blogID},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "similarity", Value: -1}}).SetLimit(int64(limit))

	cursor, err := r.similaritiesCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var similarities []models.ContentSimilarity
	if err = cursor.All(context.TODO(), &similarities); err != nil {
		return nil, err
	}

	return similarities, nil
}

func (r *recommendationMongoRepo) UpdateContentSimilarity(similarity models.ContentSimilarity) error {
	similarity.LastUpdated = time.Now()

	filter := bson.M{
		"$or": []bson.M{
			{"blog_id_1": similarity.BlogID1, "blog_id_2": similarity.BlogID2},
			{"blog_id_1": similarity.BlogID2, "blog_id_2": similarity.BlogID1},
		},
	}

	update := bson.M{"$set": similarity}
	opts := options.Update().SetUpsert(true)

	_, err := r.similaritiesCollection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

func (r *recommendationMongoRepo) GetContentSimilarities(blogID string) ([]models.ContentSimilarity, error) {
	return r.GetSimilarContent(blogID, 100)
}

// User Recommendations

func (r *recommendationMongoRepo) CreateUserRecommendation(recommendation models.UserRecommendation) error {
	recommendation.GeneratedAt = time.Now()
	recommendation.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	_, err := r.recommendationsCollection.InsertOne(context.TODO(), recommendation)
	return err
}

func (r *recommendationMongoRepo) GetUserRecommendations(userID string, limit int, category string) ([]models.UserRecommendation, error) {
	filter := bson.M{"user_id": userID, "expires_at": bson.M{"$gt": time.Now()}}

	if category != "" && category != models.CategoryAll {
		filter["category"] = category
	}

	opts := options.Find().SetSort(bson.D{{Key: "score", Value: -1}}).SetLimit(int64(limit))

	cursor, err := r.recommendationsCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var recommendations []models.UserRecommendation
	if err = cursor.All(context.TODO(), &recommendations); err != nil {
		return nil, err
	}

	return recommendations, nil
}

func (r *recommendationMongoRepo) UpdateRecommendationViewed(recommendationID string) error {
	objectID, err := primitive.ObjectIDFromHex(recommendationID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"is_viewed": true,
			"viewed_at": now,
		},
	}

	_, err = r.recommendationsCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		update,
	)
	return err
}

func (r *recommendationMongoRepo) DeleteExpiredRecommendations() error {
	filter := bson.M{"expires_at": bson.M{"$lt": time.Now()}}
	_, err := r.recommendationsCollection.DeleteMany(context.TODO(), filter)
	return err
}

func (r *recommendationMongoRepo) GetRecommendationStats(userID string) (models.RecommendationStats, error) {
	var stats models.RecommendationStats

	filter := bson.M{"user_id": userID}
	err := r.statsCollection.FindOne(context.TODO(), filter).Decode(&stats)
	if err == mongo.ErrNoDocuments {
		// Create new stats record
		stats = models.RecommendationStats{
			UserID:                 userID,
			TotalRecommendations:   0,
			ViewedRecommendations:  0,
			ClickedRecommendations: 0,
			AverageScore:           0.0,
			LastGeneratedAt:        time.Now(),
			UpdatedAt:              time.Now(),
		}
		_, err = r.statsCollection.InsertOne(context.TODO(), stats)
	}

	return stats, err
}

// User Interests

func (r *recommendationMongoRepo) UpdateUserInterest(interest models.UserInterest) error {
	interest.UpdatedAt = time.Now()

	filter := bson.M{"user_id": interest.UserID, "topic": interest.Topic}
	update := bson.M{"$set": interest}
	opts := options.Update().SetUpsert(true)

	_, err := r.interestsCollection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

func (r *recommendationMongoRepo) GetUserInterests(userID string) ([]models.UserInterest, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "weight", Value: -1}})

	cursor, err := r.interestsCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var interests []models.UserInterest
	if err = cursor.All(context.TODO(), &interests); err != nil {
		return nil, err
	}

	return interests, nil
}

func (r *recommendationMongoRepo) GetTopUserInterests(userID string, limit int) ([]models.UserInterest, error) {
	interests, err := r.GetUserInterests(userID)
	if err != nil {
		return nil, err
	}

	if limit > len(interests) {
		limit = len(interests)
	}

	return interests[:limit], nil
}

// Content Analysis

func (r *recommendationMongoRepo) GetPopularTags(limit int) ([]string, error) {
	pipeline := []bson.M{
		{"$unwind": "$tags"},
		{"$group": bson.M{
			"_id":   "$tags",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": int64(limit)},
		{"$project": bson.M{"tag": "$_id"}},
	}

	cursor, err := r.blogsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tags []string
	for cursor.Next(context.TODO()) {
		var result struct {
			Tag string `bson:"tag"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		tags = append(tags, result.Tag)
	}

	return tags, nil
}

func (r *recommendationMongoRepo) GetTrendingBlogs(limit int) ([]models.Blog, error) {
	// Get blogs with high engagement in the last 7 days
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	filter := bson.M{
		"is_published": true,
		"created_at":   bson.M{"$gte": sevenDaysAgo},
	}

	opts := options.Find().SetSort(bson.D{
		{Key: "view_count", Value: -1},
		{Key: "likes", Value: -1},
	}).SetLimit(int64(limit))

	cursor, err := r.blogsCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (r *recommendationMongoRepo) GetPopularAuthors(limit int) ([]string, error) {
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":         "$author_id",
			"total_views": bson.M{"$sum": "$view_count"},
			"total_likes": bson.M{"$sum": "$likes"},
		}},
		{"$addFields": bson.M{
			"engagement_score": bson.M{"$add": []interface{}{"$total_views", bson.M{"$multiply": []interface{}{"$total_likes", 2}}}},
		}},
		{"$sort": bson.M{"engagement_score": -1}},
		{"$limit": int64(limit)},
		{"$project": bson.M{"author_id": "$_id"}},
	}

	cursor, err := r.blogsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var authors []string
	for cursor.Next(context.TODO()) {
		var result struct {
			AuthorID string `bson:"author_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		authors = append(authors, result.AuthorID)
	}

	return authors, nil
}

// Background Processing

func (r *recommendationMongoRepo) GetBlogsForSimilarityCalculation(limit int) ([]models.Blog, error) {
	filter := bson.M{"is_published": true}
	opts := options.Find().SetLimit(int64(limit))

	cursor, err := r.blogsCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (r *recommendationMongoRepo) GetUsersForRecommendationGeneration(limit int) ([]string, error) {
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$user_id"}},
		{"$limit": int64(limit)},
		{"$project": bson.M{"user_id": "$_id"}},
	}

	cursor, err := r.behaviorsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var userIDs []string
	for cursor.Next(context.TODO()) {
		var result struct {
			UserID string `bson:"user_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		userIDs = append(userIDs, result.UserID)
	}

	return userIDs, nil
}

// Utility Methods

func (r *recommendationMongoRepo) CleanupOldBehaviors(daysOld int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)
	filter := bson.M{"created_at": bson.M{"$lt": cutoffDate}}
	_, err := r.behaviorsCollection.DeleteMany(context.TODO(), filter)
	return err
}

func (r *recommendationMongoRepo) CleanupOldSimilarities(daysOld int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)
	filter := bson.M{"last_updated": bson.M{"$lt": cutoffDate}}
	_, err := r.similaritiesCollection.DeleteMany(context.TODO(), filter)
	return err
}

// Helper methods

func (r *recommendationMongoRepo) getBlogByID(blogID string) (models.Blog, error) {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return models.Blog{}, err
	}

	var blog models.Blog
	err = r.blogsCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&blog)
	return blog, err
}

func calculateTagOverlap(tags1, tags2 []string) float64 {
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
