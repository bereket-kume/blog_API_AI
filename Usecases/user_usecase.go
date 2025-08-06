package usecases

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	repo interfaces.UserRepository
}

func NewUserUsecase(repo interfaces.UserRepository) interfaces.UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) UpdateProfile(ctx context.Context, id primitive.ObjectID, user models.User) (models.User, error) {
	return u.repo.UpdateUserProfile(ctx, id, user)
}

func (u *userUsecase) GetProfile(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	return u.repo.GetUserByID(ctx, id)
}
