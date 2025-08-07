package controllers

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"blog-api/Infrastructure/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userUsecase interfaces.UserUsecase

// This function initializes the controller with a usecase implementation
func InitUserController(u interfaces.UserUsecase) {
	userUsecase = u
}

// Handler to update user profile
func UpdateUserProfile(c *gin.Context) {
	// TODO: Replace with actual user ID from auth middleware
	userID, _ := primitive.ObjectIDFromHex("64eabf5b17c2f7e8b9cc6e9a")

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
	// TODO: Replace with actual user ID from auth middleware
	userID, _ := primitive.ObjectIDFromHex("6893b95f2f0bd5cf28b04d01")

	user, err := userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SendSuccess(c, "Profile fetched", user)
}
