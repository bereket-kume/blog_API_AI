package interfaces

import "blog-api/Domain/models"

type AISuggestionRepository interface {
	CreateAISuggestion(suggestion models.AISuggestion) (models.AISuggestion, error)
	GetAISuggestionByID(suggestionID string) (models.AISuggestion, error)
	GetAISuggestionsByUserID(userID string, page, limit int) ([]models.AISuggestion, error)
	UpdateAISuggestion(suggestion models.AISuggestion) (models.AISuggestion, error)
	DeleteAISuggestion(suggestionID string) error
	
	// Get suggestions by status
	GetAISuggestionsByStatus(userID string, status string, page, limit int) ([]models.AISuggestion, error)
	
	// Convert suggestion to draft (creates a blog draft from suggestion)
	ConvertSuggestionToDraft(suggestionID string, userID string) (models.Blog, error)
}
