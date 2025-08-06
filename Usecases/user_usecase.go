package usecases

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"errors"
	"log"
	"net/mail"
	"time"
)

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type OutPutToken struct {
	Access_token  string
	Refresh_token string
}

type UserUsecaseInterface interface {
	Register(user models.User) error
	Login(user models.User) (OutPutToken, error)
	Promote(email string) error
	Demote(email string) error
	RefreshToken(token string) (string, error)
}

type userUsecase struct {
	repo         interfaces.UserRepository
	hasher       interfaces.Hasher
	tokenService interfaces.TokenService
	tokenRepo    interfaces.TokenRepository
}

func NewUserUsecase(repo interfaces.UserRepository, hasher interfaces.Hasher, tokenService interfaces.TokenService, tokenRepo interfaces.TokenRepository) *userUsecase {
	return &userUsecase{repo: repo, hasher: hasher, tokenService: tokenService, tokenRepo: tokenRepo}
}

func (uc *userUsecase) RefreshToken(tokenStr string) (string, error) {
	// 1️⃣ Verify JWT refresh token (signature, format)
	token, err := uc.tokenService.VerifyRefreshToken(tokenStr)
	if err != nil {
		return "", err
	}

	// 2️⃣ Expiry check
	if token.ExpiresAt.Before(time.Now()) {
		if err := uc.tokenRepo.DeleteToken(token.TokenID); err != nil {
			// optional: log error, but don't return it
			log.Printf("failed to delete expired token: %v", err)
		}
		return "", errors.New("the refresh token expired")
	}

	// 3️⃣ Fetch token from DB
	dbToken, err := uc.tokenRepo.GetToken(token.TokenID)
	if err != nil {
		return "", err
	}
	if dbToken == nil {
		return "", errors.New("refresh token not found")
	}

	// 4️⃣ Verify stored hashed token matches the raw token string
	if !uc.tokenService.VerifyToken(dbToken.Token, tokenStr) {
		return "", errors.New("invalid refresh token")
	}

	// 6️⃣ Generate new access token
	accessToken, err := uc.tokenService.GenerateAccessToken(token.UserID, token.Email, token.Role)
	if err != nil {
		return "", err
	}
	return accessToken, nil

}

func (uc *userUsecase) Register(user models.User) error {
	if len(user.Username) < 3 || len(user.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if !isValidEmail(user.Email) {
		return errors.New("invalid email address")
	}

	if len(user.Password) < 8 {
		return errors.New("minimum length of password")
	}

	hashed_pw, err := uc.hasher.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed_pw
	num, err := uc.repo.CountUsers()
	if err != nil {
		return err
	}
	if num == 0 {
		user.Role = "superadmin"
	} else {
		user.Role = "user"
	}
	user.Verified = true //for time
	err = uc.repo.Insert(&user)
	return err

}

func (uc *userUsecase) Login(user models.User) (OutPutToken, error) {
	existing_user, err := uc.repo.FindByEmail(user.Email)
	if err != nil {
		return OutPutToken{}, err
	}

	if !existing_user.Verified {
		return OutPutToken{}, errors.New("user not verified")
	}

	if !uc.hasher.VerifyPassword(existing_user.Password, user.Password) {
		return OutPutToken{}, errors.New("incorrect password")
	}

	access_token, err := uc.tokenService.GenerateAccessToken(existing_user.ID, existing_user.Email, existing_user.Role)
	if err != nil {
		return OutPutToken{}, err
	}
	refresh_token, err := uc.tokenService.GenerateRefreshToken(existing_user.ID, existing_user.Email, existing_user.Role)
	if err != nil {
		return OutPutToken{}, err
	}
	refresh_tokenStr := refresh_token.Token
	refresh_token.Token = uc.tokenService.HashToken(refresh_token.Token)
	uc.tokenRepo.CreateToken(refresh_token)

	return OutPutToken{access_token, refresh_tokenStr}, err
}

func (uc *userUsecase) Promote(email string) error {
	user, err := uc.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	if user.Role == "superadmin" {
		return errors.New("superadmin cannot be promoted")
	}
	err = uc.repo.UpdateRole(email, "admin")
	return err
}

func (uc *userUsecase) Demote(email string) error {
	user, err := uc.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	if user.Role == "superadmin" {
		return errors.New("superadmin cannot be demoted")
	}
	err = uc.repo.UpdateRole(email, "user")
	return err
}
