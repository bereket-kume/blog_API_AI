package interfaces

import (
	"blog-api/Domain/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecase interface {
	UpdateProfile(ctx context.Context, id primitive.ObjectID, user models.User) (models.User, error)
	GetProfile(ctx context.Context, id primitive.ObjectID) (models.User, error)
}
