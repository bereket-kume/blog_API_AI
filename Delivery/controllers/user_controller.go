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

func InitUserController(u interfaces.UserUsecase) {
	userUsecase = u
}

func UpdateUserProfile(c *gin.Context) {
	// mock user ID until auth is ready
	userID, _ := primitive.ObjectIDFromHex("6893b95f2f0bd5cf28b04d01") // replace with test user ID

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

	utils.SendSuccess(c, http.StatusOK, "Profile updated", updated)
}

func GetUserProfile(c *gin.Context) {
	userID, _ := primitive.ObjectIDFromHex("6893b95f2f0bd5cf28b04d01") // replace with test user ID

	user, err := userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Profile fetched", user)
}
