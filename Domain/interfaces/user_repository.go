package interfaces

import (
	"context"

	"blog-api/Domain/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	// From version 1
	UpdateUserProfile(ctx context.Context, id primitive.ObjectID, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (models.User, error)

	// From version 2
	Insert(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	UpdatePass(email string, passwordHash string) error
	UpdateRole(email, role string) error
	Delete(email string) error
	Verify(email string) error
	CountUsers() (int64, error)
}
