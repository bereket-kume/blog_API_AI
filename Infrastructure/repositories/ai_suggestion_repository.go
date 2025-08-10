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

type aiSuggestionMongoRepo struct {
	collection *mongo.Collection
}

func NewAISuggestionMongoRepo(col *mongo.Collection) *aiSuggestionMongoRepo {
	return &aiSuggestionMongoRepo{collection: col}
}

// CreateAISuggestion creates a new AI suggestion
func (ar *aiSuggestionMongoRepo) CreateAISuggestion(suggestion models.AISuggestion) (models.AISuggestion, error) {
	// Generate a new ObjectID for MongoDB
	objectID := primitive.NewObjectID()
	suggestion.ID = objectID.Hex()
	suggestion.CreatedAt = time.Now()
	suggestion.UpdatedAt = time.Now()

	// Convert to MongoDB model for insertion
	suggestionModel := bson.M{
		"_id":               objectID,
		"user_id":           suggestion.UserID,
		"input_topic":       suggestion.InputTopic,
		"keywords":          suggestion.Keywords,
		"tone":              suggestion.Tone,
		"suggested_content": suggestion.SuggestedContent,
		"suggestions":       suggestion.Suggestions,
		"status":            suggestion.Status,
		"created_at":        suggestion.CreatedAt,
		"updated_at":        suggestion.UpdatedAt,
	}

	_, err := ar.collection.InsertOne(context.TODO(), suggestionModel)
	if err != nil {
		return models.AISuggestion{}, err
	}

	return suggestion, nil
}

// GetAISuggestionByID retrieves an AI suggestion by its ID
func (ar *aiSuggestionMongoRepo) GetAISuggestionByID(suggestionID string) (models.AISuggestion, error) {
	objectID, err := primitive.ObjectIDFromHex(suggestionID)
	if err != nil {
		return models.AISuggestion{}, err
	}

	var suggestion models.AISuggestion
	err = ar.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&suggestion)
	if err != nil {
		return models.AISuggestion{}, err
	}

	return suggestion, nil
}

// GetAISuggestionsByUserID retrieves AI suggestions for a specific user with pagination
func (ar *aiSuggestionMongoRepo) GetAISuggestionsByUserID(userID string, page, limit int) ([]models.AISuggestion, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"user_id": userID}
	cursor, err := ar.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var suggestions []models.AISuggestion
	if err = cursor.All(context.TODO(), &suggestions); err != nil {
		return nil, err
	}

	return suggestions, nil
}

// UpdateAISuggestion updates an existing AI suggestion
func (ar *aiSuggestionMongoRepo) UpdateAISuggestion(suggestion models.AISuggestion) (models.AISuggestion, error) {
	suggestion.UpdatedAt = time.Now()

	objectID, err := primitive.ObjectIDFromHex(suggestion.ID)
	if err != nil {
		return models.AISuggestion{}, err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"user_id":           suggestion.UserID,
		"input_topic":       suggestion.InputTopic,
		"keywords":          suggestion.Keywords,
		"tone":              suggestion.Tone,
		"suggested_content": suggestion.SuggestedContent,
		"suggestions":       suggestion.Suggestions,
		"status":            suggestion.Status,
		"updated_at":        suggestion.UpdatedAt,
	}}

	_, err = ar.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return models.AISuggestion{}, err
	}

	return suggestion, nil
}

// DeleteAISuggestion deletes an AI suggestion
func (ar *aiSuggestionMongoRepo) DeleteAISuggestion(suggestionID string) error {
	objectID, err := primitive.ObjectIDFromHex(suggestionID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = ar.collection.DeleteOne(context.TODO(), filter)
	return err
}

// GetAISuggestionsByStatus retrieves AI suggestions by status for a specific user
func (ar *aiSuggestionMongoRepo) GetAISuggestionsByStatus(userID string, status string, page, limit int) ([]models.AISuggestion, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{
		"user_id": userID,
		"status":  status,
	}

	cursor, err := ar.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var suggestions []models.AISuggestion
	if err = cursor.All(context.TODO(), &suggestions); err != nil {
		return nil, err
	}

	return suggestions, nil
}

// ConvertSuggestionToDraft converts an AI suggestion to a blog draft
func (ar *aiSuggestionMongoRepo) ConvertSuggestionToDraft(suggestionID string, userID string) (models.Blog, error) {
	// This method will be implemented in the use case layer
	// as it requires both AI suggestion and blog repository access
	return models.Blog{}, nil
}
