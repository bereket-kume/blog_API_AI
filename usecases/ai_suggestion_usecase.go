package usecases

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"errors"
	"time"
)

type aiSuggestionUseCase struct {
	aiSuggestionRepo interfaces.AISuggestionRepository
	blogRepo         interfaces.BlogRepository
}

func NewAISuggestionUseCase(aiSuggestionRepo interfaces.AISuggestionRepository, blogRepo interfaces.BlogRepository) interfaces.AISuggestionUseCase {
	return &aiSuggestionUseCase{
		aiSuggestionRepo: aiSuggestionRepo,
		blogRepo:         blogRepo,
	}
}

func (a *aiSuggestionUseCase) CreateAISuggestion(suggestion models.AISuggestion) (models.AISuggestion, error) {
	// Set timestamps
	now := time.Now()
	suggestion.CreatedAt = now
	suggestion.UpdatedAt = now
	
	return a.aiSuggestionRepo.CreateAISuggestion(suggestion)
}

func (a *aiSuggestionUseCase) GetAISuggestionByID(suggestionID string, userID string) (models.AISuggestion, error) {
	suggestion, err := a.aiSuggestionRepo.GetAISuggestionByID(suggestionID)
	if err != nil {
		return models.AISuggestion{}, err
	}
	
	// Ensure user can only access their own suggestions
	if suggestion.UserID != userID {
		return models.AISuggestion{}, errors.New("unauthorized access to suggestion")
	}
	
	return suggestion, nil
}

func (a *aiSuggestionUseCase) GetAISuggestionsByUserID(userID string, page, limit int) ([]models.AISuggestion, error) {
	return a.aiSuggestionRepo.GetAISuggestionsByUserID(userID, page, limit)
}

func (a *aiSuggestionUseCase) UpdateAISuggestion(suggestion models.AISuggestion, userID string) (models.AISuggestion, error) {
	// Get existing suggestion to verify ownership
	existing, err := a.aiSuggestionRepo.GetAISuggestionByID(suggestion.ID)
	if err != nil {
		return models.AISuggestion{}, err
	}
	
	// Ensure user can only update their own suggestions
	if existing.UserID != userID {
		return models.AISuggestion{}, errors.New("unauthorized access to suggestion")
	}
	
	// Update timestamp
	suggestion.UpdatedAt = time.Now()
	suggestion.UserID = userID // Ensure userID doesn't change
	
	return a.aiSuggestionRepo.UpdateAISuggestion(suggestion)
}

func (a *aiSuggestionUseCase) DeleteAISuggestion(suggestionID string, userID string) error {
	// Get existing suggestion to verify ownership
	existing, err := a.aiSuggestionRepo.GetAISuggestionByID(suggestionID)
	if err != nil {
		return err
	}
	
	// Ensure user can only delete their own suggestions
	if existing.UserID != userID {
		return errors.New("unauthorized access to suggestion")
	}
	
	return a.aiSuggestionRepo.DeleteAISuggestion(suggestionID)
}

func (a *aiSuggestionUseCase) GetAISuggestionsByStatus(userID string, status string, page, limit int) ([]models.AISuggestion, error) {
	return a.aiSuggestionRepo.GetAISuggestionsByStatus(userID, status, page, limit)
}

func (a *aiSuggestionUseCase) ConvertSuggestionToDraft(suggestionID string, userID string) (models.Blog, error) {
	// Get the suggestion
	suggestion, err := a.GetAISuggestionByID(suggestionID, userID)
	if err != nil {
		return models.Blog{}, err
	}
	
	// Create a blog draft from the suggestion
	blog := models.Blog{
		Title:       suggestion.SuggestedContent, // Use the first suggestion as title
		Content:     "", // Empty content for draft
		AuthorID:    userID,
		AuthorName:  "", // Will be set by controller
		Tags:        suggestion.Keywords,
		IsPublished: false, // Draft
		ViewCount:   0,
		Likes:       0,
		Dislikes:    0,
		Comments:    []models.Comment{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// Create the blog draft
	createdBlog, err := a.blogRepo.CreateBlog(blog)
	if err != nil {
		return models.Blog{}, err
	}
	
	// Update suggestion status to converted
	suggestion.Status = models.AISuggestionStatusConvertedToDraft
	suggestion.UpdatedAt = time.Now()
	_, err = a.aiSuggestionRepo.UpdateAISuggestion(suggestion)
	if err != nil {
		return models.Blog{}, err
	}
	
	return createdBlog, nil
}

func (a *aiSuggestionUseCase) SaveAISuggestion(userID string, inputTopic string, keywords []string, tone string, suggestions []string) (models.AISuggestion, error) {
	// Create a new AI suggestion
	suggestion := models.AISuggestion{
		UserID:           userID,
		InputTopic:       inputTopic,
		Keywords:         keywords,
		Tone:             tone,
		SuggestedContent: "", // Will be populated from suggestions
		Suggestions:      suggestions,
		Status:           models.AISuggestionStatusSaved,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	
	// If there are suggestions, use the first one as suggested content
	if len(suggestions) > 0 {
		suggestion.SuggestedContent = suggestions[0]
	}
	
	return a.aiSuggestionRepo.CreateAISuggestion(suggestion)
}
