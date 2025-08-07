// infrastructure/mongo/user_model.go
package db_models

import (
	"blog-api/Domain/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserModel represents the MongoDB-specific user model
// This handles the conversion between domain and infrastructure concerns
type UserModel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
	Role         string             `bson:"role"`
	Verified     bool               `bson:"verified"`
}

// FromDomainUser converts a domain User to a MongoDB UserModel
func FromDomainUser(u *models.User) *UserModel {
	var objectID primitive.ObjectID
	if u.ID != "" {
		if id, err := primitive.ObjectIDFromHex(u.ID); err == nil {
			objectID = id
		}
	}

	return &UserModel{
		ID:           objectID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.Password,
		Role:         u.Role,
		Verified:     u.Verified,
	}
}

// ToDomainUser converts a MongoDB UserModel to a domain User
func ToDomainUser(m *UserModel) *models.User {
	return &models.User{
		ID:       m.ID.Hex(),
		Username: m.Username,
		Email:    m.Email,
		Password: m.PasswordHash,
		Role:     m.Role,
		Verified: m.Verified,
	}
}
