package usecases

import (
	"blog-api/Domain/models"
	"context"
)

func (u *userUsecase) UpdateProfile(ctx context.Context, id string, user models.User) (models.User, error) {
	return u.repo.UpdateUserProfile(ctx, id, user)
}

func (u *userUsecase) GetProfile(ctx context.Context, id string) (models.User, error) {
	return u.repo.GetUserByID(ctx, id)
}
