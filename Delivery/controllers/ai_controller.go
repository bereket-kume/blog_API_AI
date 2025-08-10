package controllers

import (
	"blog-api/Domain/interfaces"
	"blog-api/Infrastructure/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AIRequest struct {
	Keywords []string `json:"keywords"`
	Tone     string   `json:"tone"`
	Content  string   `json:"content"`
}

type AISuggestionRequest struct {
	BlogContent string   `json:"blog_content"`
	Keywords    []string `json:"keywords" binding:"required"`
	Tone        string   `json:"tone"`
	Type        string   `json:"type"` // "improvement", "ideas", "title", "summary"
}

type AISuggestionResponse struct {
	Suggestions []string `json:"suggestions"`
	Type        string   `json:"type"`
	Message     string   `json:"message"`
}

type AISuggestionController struct {
	aiSuggestionUC interfaces.AISuggestionUseCase
}

func NewAISuggestionController(aiSuggestionUC interfaces.AISuggestionUseCase) *AISuggestionController {
	return &AISuggestionController{aiSuggestionUC: aiSuggestionUC}
}

// GenerateAISuggestion handles AI content suggestions for blog posts
func GenerateAISuggestion(c *gin.Context) {
	var req AISuggestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid input format: "+err.Error())
		return
	}

	// Validate request
	if len(req.Keywords) == 0 {
		utils.SendError(c, http.StatusBadRequest, "Keywords are required")
		return
	}

	// Set default tone if not provided
	if req.Tone == "" {
		req.Tone = "professional"
	}

	// Set default type if not provided
	if req.Type == "" {
		req.Type = "improvement"
	}

	// Prepare payload for AI service (matching the AI service interface)
	aiPayload := map[string]interface{}{
		"keywords": req.Keywords,
		"tone":     req.Tone,
	}

	// Marshal the input to JSON for sending to the AI microservice
	payloadBytes, err := json.Marshal(aiPayload)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to encode request")
		return
	}

	// Get AI service URL from environment
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		utils.SendError(c, http.StatusInternalServerError, "AI_SERVICE_URL not configured")
		return
	}

	// Make request to AI service
	resp, err := http.Post(aiURL+"/generate", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "AI service request failed: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Check if AI service responded successfully
	if resp.StatusCode != http.StatusOK {
		utils.SendError(c, http.StatusInternalServerError, "AI service returned error status: "+resp.Status)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to parse AI response")
		return
	}

	// Create response
	aiResponse := AISuggestionResponse{
		Suggestions: []string{},
		Type:        req.Type,
		Message:     "AI suggestions generated successfully",
	}

	// Extract data from AI service response (matching the AI service response structure)
	if title, ok := result["title"].(string); ok && title != "" {
		aiResponse.Suggestions = append(aiResponse.Suggestions, "Title: "+title)
	}

	if audience, ok := result["audience"].(string); ok && audience != "" {
		aiResponse.Suggestions = append(aiResponse.Suggestions, "Target Audience: "+audience)
	}

	if headlines, ok := result["headlines"].([]interface{}); ok {
		for _, headline := range headlines {
			if str, ok := headline.(string); ok {
				aiResponse.Suggestions = append(aiResponse.Suggestions, "• "+str)
			}
		}
	}

	utils.SendSuccess(c, aiResponse.Message, aiResponse)
}

// SaveAISuggestion saves the AI suggestion to the database
func (ctrl *AISuggestionController) SaveAISuggestion(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		InputTopic  string   `json:"input_topic" binding:"required"`
		Keywords    []string `json:"keywords" binding:"required"`
		Tone        string   `json:"tone"`
		Suggestions []string `json:"suggestions" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid input format: "+err.Error())
		return
	}

	// Set default tone if not provided
	if req.Tone == "" {
		req.Tone = "professional"
	}

	// Save the suggestion to database
	suggestion, err := ctrl.aiSuggestionUC.SaveAISuggestion(
		userID.(string),
		req.InputTopic,
		req.Keywords,
		req.Tone,
		req.Suggestions,
	)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to save suggestion: "+err.Error())
		return
	}

	utils.SendSuccess(c, "AI suggestion saved successfully", suggestion)
}

// GetAISuggestions retrieves AI suggestions for the authenticated user
func (ctrl *AISuggestionController) GetAISuggestions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	suggestions, err := ctrl.aiSuggestionUC.GetAISuggestionsByUserID(userID.(string), page, limit)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve suggestions: "+err.Error())
		return
	}

	utils.SendSuccess(c, "AI suggestions retrieved successfully", suggestions)
}

// GetAISuggestionsByStatus retrieves AI suggestions by status
func (ctrl *AISuggestionController) GetAISuggestionsByStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	status := c.Param("status")
	if status == "" {
		utils.SendError(c, http.StatusBadRequest, "Status parameter is required")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	suggestions, err := ctrl.aiSuggestionUC.GetAISuggestionsByStatus(userID.(string), status, page, limit)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve suggestions: "+err.Error())
		return
	}

	utils.SendSuccess(c, "AI suggestions retrieved successfully", suggestions)
}

// ConvertSuggestionToDraft converts an AI suggestion to a blog draft
func (ctrl *AISuggestionController) ConvertSuggestionToDraft(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	suggestionID := c.Param("id")
	if suggestionID == "" {
		utils.SendError(c, http.StatusBadRequest, "Suggestion ID is required")
		return
	}

	// Convert suggestion to draft
	blog, err := ctrl.aiSuggestionUC.ConvertSuggestionToDraft(suggestionID, userID.(string))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to convert suggestion to draft: "+err.Error())
		return
	}

	utils.SendSuccess(c, "Suggestion converted to draft successfully", blog)
}

// DeleteAISuggestion deletes an AI suggestion
func (ctrl *AISuggestionController) DeleteAISuggestion(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	suggestionID := c.Param("id")
	if suggestionID == "" {
		utils.SendError(c, http.StatusBadRequest, "Suggestion ID is required")
		return
	}

	err := ctrl.aiSuggestionUC.DeleteAISuggestion(suggestionID, userID.(string))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to delete suggestion: "+err.Error())
		return
	}

	utils.SendSuccess(c, "AI suggestion deleted successfully", nil)
}

// GenerateContentIdeas handles AI-generated content ideas
func GenerateContentIdeas(c *gin.Context) {
	var req struct {
		Keywords []string `json:"keywords"`
		Tone     string   `json:"tone"`
		Category string   `json:"category"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid input format: "+err.Error())
		return
	}

	// Set default tone if not provided
	if req.Tone == "" {
		req.Tone = "professional"
	}

	// Marshal the input to JSON for sending to the AI microservice
	payloadBytes, err := json.Marshal(map[string]interface{}{
		"keywords": req.Keywords,
		"tone":     req.Tone,
	})

	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to encode request")
		return
	}

	// Get AI service URL from environment
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		utils.SendError(c, http.StatusInternalServerError, "AI_SERVICE_URL not configured")
		return
	}

	// Make request to AI service
	resp, err := http.Post(aiURL+"/generate", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "AI service request failed: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Check if AI service responded successfully
	if resp.StatusCode != http.StatusOK {
		utils.SendError(c, http.StatusInternalServerError, "AI service returned error status: "+resp.Status)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to parse AI response")
		return
	}

	// Create response
	aiResponse := AISuggestionResponse{
		Suggestions: []string{},
		Type:        "ideas",
		Message:     "Content ideas generated successfully",
	}

	// Extract data from AI service response (matching the AI service response structure)
	if title, ok := result["title"].(string); ok && title != "" {
		aiResponse.Suggestions = append(aiResponse.Suggestions, "Title: "+title)
	}

	if audience, ok := result["audience"].(string); ok && audience != "" {
		aiResponse.Suggestions = append(aiResponse.Suggestions, "Target Audience: "+audience)
	}

	if headlines, ok := result["headlines"].([]interface{}); ok {
		for _, headline := range headlines {
			if str, ok := headline.(string); ok {
				aiResponse.Suggestions = append(aiResponse.Suggestions, "• "+str)
			}
		}
	}

	utils.SendSuccess(c, aiResponse.Message, aiResponse)
}
