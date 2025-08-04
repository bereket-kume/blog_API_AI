package controllers

import (
	"blog-api/Domain/models"
	Database "blog-api/Infrastructure/database"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserProfile(c *gin.Context) {
	email := c.MustGet("email").(string)

	collection := Database.GetUserCollection()
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUserProfile(c *gin.Context) {
	email := c.MustGet("email").(string)

	var updateData struct {
		Bio        string `json:"bio"`
		ProfilePic string `json:"profile_pic"`
		Contact    string `json:"contact"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	collection := Database.GetUserCollection()
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"bio":         updateData.Bio,
			"profile_pic": updateData.ProfilePic,
			"contact":     updateData.Contact,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}
