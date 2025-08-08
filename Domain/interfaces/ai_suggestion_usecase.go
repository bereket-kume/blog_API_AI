package interfaces

import "blog-api/Domain/models"

type AISuggestionUseCase interface {
	CreateAISuggestion(suggestion models.AISuggestion) (models.AISuggestion, error)
	GetAISuggestionByID(suggestionID string, userID string) (models.AISuggestion, error)
	GetAISuggestionsByUserID(userID string, page, limit int) ([]models.AISuggestion, error)
	UpdateAISuggestion(suggestion models.AISuggestion, userID string) (models.AISuggestion, error)
	DeleteAISuggestion(suggestionID string, userID string) error
	
	// Get suggestions by status
	GetAISuggestionsByStatus(userID string, status string, page, limit int) ([]models.AISuggestion, error)
	
	// Convert suggestion to draft
	ConvertSuggestionToDraft(suggestionID string, userID string) (models.Blog, error)
	
	// Save AI suggestion with status
	SaveAISuggestion(userID string, inputTopic string, keywords []string, tone string, suggestions []string) (models.AISuggestion, error)
}
