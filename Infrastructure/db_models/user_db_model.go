// infrastructure/mongo/user_model.go
package db_models

import (
	"blog-api/Domain/models"
)

type UserModel struct {
	ID           string `bson:"_id,omitempty"`
	Username     string `bson:"username"`
	Email        string `bson:"email"`
	PasswordHash string `bson:"password_hash"`
	Role         string `bson:"role"`
	Verified     bool   `bson:"verified"`
}

func FromDomainUser(u *models.User) *UserModel {
	return &UserModel{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.Password,
		Role:         u.Role,
		Verified:     u.Verified,
	}
}

func ToDomainUser(m *UserModel) *models.User {
	return &models.User{
		ID:       m.ID,
		Username: m.Username,
		Email:    m.Email,
		Password: m.PasswordHash,
		Role:     m.Role,
		Verified: m.Verified,
	}
}
