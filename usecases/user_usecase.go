package usecases

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
	"context"
	"errors"
	"fmt"
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
	Logout(refresh_token string) error
	VerifyEmail(tokenStr string) error
	RequestPasswordReset(email string) error
	ResetPassword(resetToken, newPassword string) error
	UpdateProfile(ctx context.Context, id string, user models.User) (models.User, error)
	GetProfile(ctx context.Context, id string) (models.User, error)
}

type userUsecase struct {
	repo         interfaces.UserRepository
	hasher       interfaces.Hasher
	tokenService interfaces.TokenService
	tokenRepo    interfaces.TokenRepository
	emailService interfaces.EmailService
}

func NewUserUsecase(repo interfaces.UserRepository, hasher interfaces.Hasher, tokenService interfaces.TokenService, tokenRepo interfaces.TokenRepository, emailService interfaces.EmailService) *userUsecase {
	return &userUsecase{repo: repo, hasher: hasher, tokenService: tokenService, tokenRepo: tokenRepo, emailService: emailService}
}

func (uc *userUsecase) Logout(refreshToken string) error {
	// 1️⃣ Verify and parse refresh token
	token, err := uc.tokenService.VerifyRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	// 2️⃣ Expiry check
	if token.ExpiresAt.Before(time.Now()) {
		_ = uc.tokenRepo.DeleteToken(token.TokenID) // best effort
		return errors.New("the refresh token expired")
	}

	// 3️⃣ Fetch from DB
	dbToken, err := uc.tokenRepo.GetToken(token.TokenID)
	if err != nil {
		return err
	}
	if dbToken == nil {
		return errors.New("refresh token not found")
	}

	// 4️⃣ Verify token hash matches
	if !uc.tokenService.VerifyToken(dbToken.Token, refreshToken) {
		return errors.New("invalid refresh token")
	}

	// 5️⃣ Delete token from DB (logout)
	return uc.tokenRepo.DeleteToken(dbToken.ID)
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
	// 1. Validate email format
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
	_, err = uc.repo.FindByEmail(user.Email)
	if err == nil {
		return errors.New("email already exists")
	}
	user.Verified = false

	// 4. Generate verification token (UUID or JWT)
	Exp := time.Hour * 1
	verificationToken, err := uc.tokenService.GenerateRandomJWT(Exp)
	if err != nil {
		return err
	}
	verificationToken.Email = user.Email
	tokenStr := verificationToken.Token
	verificationToken.Token = uc.tokenService.HashToken(tokenStr)

	// 5. Insert user into database first
	if err := uc.repo.Insert(&user); err != nil {
		return err
	}

	// 6. Store verification token
	if err := uc.tokenRepo.CreateToken(verificationToken); err != nil {
		// If token storage fails, we should clean up the user
		log.Printf("Failed to store verification token for user %s, cleaning up: %v", user.Email, err)
		// Note: In production, you might want to implement a cleanup mechanism
		return fmt.Errorf("failed to create verification token: %w", err)
	}

	// 7. Send verification email (non-blocking for registration success)
	go func() {
		if err := uc.emailService.SendVerificationEmail(user.Username, user.Email, tokenStr); err != nil {
			log.Printf("Failed to send verification email to %s: %v", user.Email, err)
			// Don't fail registration if email fails, just log it
		}
	}()

	return nil

}

func (uc *userUsecase) Login(user models.User) (OutPutToken, error) {
	existing_user, err := uc.repo.FindByEmail(user.Email)
	if err != nil {
		return OutPutToken{}, errors.New("user not found")
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

func (uc *userUsecase) VerifyEmail(tokenStr string) error {
	// 1. Fetch token from DB
	token, err := uc.tokenService.VerifyJWT(tokenStr)
	if err != nil {
		return errors.New("invalid token")
	}
	dbToken, err := uc.tokenRepo.GetToken(token.TokenID)
	if err != nil || dbToken == nil {
		return errors.New("invalid or expired token")
	}

	// 2. Check expiration
	if dbToken.ExpiresAt.Before(time.Now()) {
		return errors.New("verification token expired")
	}

	// 3. Mark user as verified
	log.Println(dbToken.Email)
	err = uc.repo.Verify(dbToken.Email)
	return err
}

func (uc *userUsecase) RequestPasswordReset(email string) error {
	// 1️⃣ Find the user
	user, err := uc.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	// 2️⃣ Generate token (reuse tokenService)
	exp := time.Minute * 15
	resetToken, err := uc.tokenService.GenerateRandomJWT(exp)
	if err != nil {
		return err
	}
	tokenStr := resetToken.Token
	resetToken.Token = uc.tokenService.HashToken(tokenStr)
	resetToken.Email = user.Email

	// 3️⃣ Store in DB as a token entity
	if err := uc.tokenRepo.CreateToken(resetToken); err != nil {
		return err
	}

	return uc.emailService.SendPasswordResetEmail(user.Username, user.Email, tokenStr)

}

func (uc *userUsecase) ResetPassword(resetToken, newPassword string) error {
	// 1️⃣ Fetch token claims
	// (Here we just find token from DB directly since it's plain string from email)
	token, err := uc.tokenService.VerifyJWT(resetToken) // Or custom GetByTokenID
	if err != nil {
		return err
	}

	db_token, err := uc.tokenRepo.GetToken(token.TokenID)
	if err != nil {
		return err
	}

	// 2️⃣ Check expiry
	if db_token.ExpiresAt.Before(time.Now()) {
		_ = uc.tokenRepo.DeleteToken(db_token.ID)
		return errors.New("token expired")
	}

	// 3️⃣ Update password
	hashedPass, err := uc.hasher.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := uc.repo.UpdatePass(db_token.Email, hashedPass); err != nil {
		return err
	}

	// 4️⃣ Delete token after use
	return uc.tokenRepo.DeleteToken(db_token.ID)
}
