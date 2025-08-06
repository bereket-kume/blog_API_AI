package controllers

import (
	"blog-api/Infrastructure/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AIRequest struct {
	Keywords []string `json:"keywords"`
	Tone     string   `json:"tone"`
}

func GenerateAISuggestion(c *gin.Context) {
	var req AIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	// Marshal the input to JSON for sending to the AI microservice
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to encode request")
		return
	}

	// Get AI service URL from environment
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		utils.RespondWithError(c, http.StatusInternalServerError, "AI_SERVICE_URL not configured")
		return
	}

	resp, err := http.Post(aiURL+"/generate", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "AI service request failed")
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to parse AI response")
		return
	}

	utils.RespondWithSuccess(c, "AI suggestion generated successfully", result)
}
