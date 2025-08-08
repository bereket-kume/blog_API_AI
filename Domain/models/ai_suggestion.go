package models

import (
	"time"
)

// AISuggestion represents an AI-generated suggestion stored in the database
type AISuggestion struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	UserID          string    `json:"user_id" bson:"user_id"`
	InputTopic      string    `json:"input_topic" bson:"input_topic"`
	Keywords        []string  `json:"keywords" bson:"keywords"`
	Tone            string    `json:"tone" bson:"tone"`
	SuggestedContent string   `json:"suggested_content" bson:"suggested_content"`
	Suggestions     []string  `json:"suggestions" bson:"suggestions"`
	Status          string    `json:"status" bson:"status"` // "saved", "converted-to-draft", "discarded"
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// AISuggestionStatus represents the possible statuses for AI suggestions
const (
	AISuggestionStatusSaved           = "saved"
	AISuggestionStatusConvertedToDraft = "converted-to-draft"
	AISuggestionStatusDiscarded       = "discarded"
)
