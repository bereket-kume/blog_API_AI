package usecases

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"errors"
	"net/mail"
)

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type UserUsecaseInterface interface {
	Register(user models.User) error
	Login(user models.User) (string, error)
	Promote(email string) error
	Demote(email string) error
}

type userUsecase struct {
	repo         interfaces.UserRepository
	hasher       models.Hasher
	tokenService models.TokenService
}

func NewUserUsecase(repo interfaces.UserRepository, hasher models.Hasher, tokenService models.TokenService) *userUsacase {
	return &userUsecase{repo: repo, hasher: hasher, tokenService: tokenService}
}

func (uc *userUsecase) Register(user models.User) error {
	if len(user.Username) < 3 || len(user.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if !isValidEmail(user.Email) {
		return errors.New("Invalid Email address")
	}

	if len(user.Password) < 8 {
		return errors.New("Minimum length of password")
	}

	hashed_pw := uc.hasher.HashPassword(user.Password)
	user.Password = hashed_pw
	num, err := uc.repo.CountUsers()
	if num == 0 {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}
	user.Promoted_by = "null"
	err = uc.repo.Insert(user)
	return err

}

func (uc *userUsecase) Login(user models.User) (string, error) {
	existing_user, err := uc.repo.FindByEmail(user.Email)
	if err != nil {
		return "", err
	}

	if !existing_user.Verified {
		return "", errors.New("user not verified")
	}

	if !uc.hasher.VerifyPassword(existing_user.Password, user.Password) {
		return "", errors.New("incorrect password")
	}

	token, err := uc.tokenService.GenerateToken(existing_user.ID, existing_user.Email, existing_user.Role)

	return token, err
}

func (uc *userUsecase) Promote(email string) error {
	if _, err := uc.repo.FindByEmail(email); err != nil {
		return errors.New("user not found")
	}
	err := uc.repo.UpdateRole(email, "admin")
	return err
}

func (uc *userUsecase) Demote(email string) error {
	if _, err := uc.repo.FindByEmail(email); err != nil {
		return errors.New("user not found")
	}
	err := uc.repo.UpdateRole(email, "user")
	return err
}
