package interfaces

import (
	"context"

	"blog-api/Domain/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	UpdateUserProfile(ctx context.Context, id primitive.ObjectID, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (models.User, error)
}
