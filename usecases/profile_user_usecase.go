package usecases

import (
	"blog-api/Domain/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *userUsecase) UpdateProfile(ctx context.Context, id primitive.ObjectID, user models.User) (models.User, error) {
	return u.repo.UpdateUserProfile(ctx, id, user)
}

func (u *userUsecase) GetProfile(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	return u.repo.GetUserByID(ctx, id)
}
