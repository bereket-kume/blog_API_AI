package interfaces

import (
	"blog-api/Domain/models"
	"context"
)

type UserUsecase interface {
	UpdateProfile(ctx context.Context, id string, user models.User) (models.User, error)
	GetProfile(ctx context.Context, id string) (models.User, error)
}
