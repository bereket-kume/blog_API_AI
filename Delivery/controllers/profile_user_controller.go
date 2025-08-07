package controllers

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"blog-api/Infrastructure/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

var userUsecase interfaces.UserUsecase

// This function initializes the controller with a usecase implementation
func InitUserController(u interfaces.UserUsecase) {
	userUsecase = u
}

// Handler to update user profile
func UpdateUserProfile(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found")
		return
	}

	userID, ok := userIDRaw.(primitive.ObjectID)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "Invalid user ID format")
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	updated, err := userUsecase.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Update failed")
		return
	}

	utils.SendSuccess(c, "Profile updated", updated)
}

// Handler to fetch user profile
func GetUserProfile(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found")
		return
	}

	userID, ok := userIDRaw.(primitive.ObjectID)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "Invalid user ID format")
		return
	}

	user, err := userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SendSuccess(c, "Profile fetched", user)
}
