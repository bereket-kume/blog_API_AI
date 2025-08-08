package usecases

import (
	"blog-api/Domain/models"
	"blog-api/mocks" // generated mocks
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (*mocks.UserRepository, *mocks.MockHasher, *mocks.MockTokenService, *mocks.MockTokenRepository, *mocks.MockEmailService, UserUsecaseInterface) {
	repo := new(mocks.UserRepository)
	hasher := new(mocks.MockHasher)
	tokenSvc := new(mocks.MockTokenService)
	tokenRepo := new(mocks.MockTokenRepository)
	emailService := new(mocks.MockEmailService)

	uc := NewUserUsecase(repo, hasher, tokenSvc, tokenRepo, emailService)
	return repo, hasher, tokenSvc, tokenRepo, emailService, uc
}

func TestRegister_EmailSuccess_UserSaved(t *testing.T) {
	repo, hasher, tokenSvc, tokenRepo, emailService, uc := setup()

	user := models.User{
		ID:       "user-id-123",
		Username: "jane",
		Email:    "jane@example.com",
		Password: "pass1234",
		Role:     "user",
		Verified: false,
	}
	emailService.On("SendVerificationEmail", user.Username, user.Email, "plainToken").Return(nil)
	hasher.On("HashPassword", "pass1234").Return("hashed_pass", nil)
	repo.On("CountUsers", mock.Anything).Return(0, nil)
	repo.On("FindByEmail", "jane@example.com").Return(&models.User{}, errors.New("no document"))
	tokenSvc.On("GenerateRandomJWT", time.Hour*1).Return(&models.Token{
		Token: "plainToken",
	}, nil)
	tokenSvc.On("HashToken", "plainToken").Return("hashedToken")
	repo.On("Insert", mock.MatchedBy(func(u *models.User) bool {
		return u.Email == "jane@example.com" &&
			u.Username == "jane" &&
			u.Password == "hashed_pass" &&
			u.Role == "superadmin" &&
			!u.Verified
	})).Return(nil)

	tokenRepo.On("CreateToken", mock.MatchedBy(func(t *models.Token) bool {
		return t.Token == "hashedToken" && t.Email == "jane@example.com"
	})).Return(nil)

	err := uc.Register(user)

	assert.NoError(t, err)
	emailService.AssertExpectations(t)
}

func TestLogin_UserNotVerified(t *testing.T) {
	repo, _, _, _, _, uc := setup()

	existing := models.User{
		Email:    "john@example.com",
		Verified: false,
	}

	repo.On("FindByEmail", "john@example.com").Return(existing, nil)

	_, err := uc.Login(models.User{
		Email:    "john@example.com",
		Password: "password123",
	})

	assert.EqualError(t, err, "user not verified")
}
func TestLogin_success(t *testing.T) {
	repo, mockhasher, token_service, tokenRepo, _, uc := setup()

	existing := models.User{
		ID:       "123456789",
		Email:    "john@example.com",
		Verified: true,
		Password: "hashed_password",
		Role:     "user",
	}

	repo.On("FindByEmail", "john@example.com").Return(existing, nil)
	mockhasher.On("VerifyPassword", existing.Password, "password123").Return(true)
	token_service.On("GenerateAccessToken", "123456789", "john@example.com", "user").Return("access_token", nil)
	token_service.On("GenerateRefreshToken", "123456789", "john@example.com", "user").Return(&models.Token{Token: "refresh_token"}, nil)
	token_service.On("HashToken", "refresh_token").Return("hashed_refresh_token")
	tokenRepo.On("CreateToken", &models.Token{Token: "hashed_refresh_token"}).Return(nil)

	tokens, err := uc.Login(models.User{
		Email:    "john@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.Equal(t, "access_token", tokens.Access_token)
	assert.Equal(t, "refresh_token", tokens.Refresh_token)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo, _, _, _, _, uc := setup()

	repo.On("FindByEmail", "john@example.com").Return(models.User{}, errors.New("User not found"))

	_, err := uc.Login(models.User{
		Email:    "john@example.com",
		Password: "password123",
	})

	assert.EqualError(t, err, "user not found")
}
func TestLogin_PasswordIncorrect(t *testing.T) {
	repo, hasher, _, _, _, uc := setup()

	existing := models.User{
		ID:       "123456789",
		Email:    "john@example.com",
		Verified: true,
		Password: "hashed_password",
		Role:     "user",
	}

	repo.On("FindByEmail", "john@example.com").Return(existing, nil)
	hasher.On("VerifyPassword", existing.Password, "password123").Return(false)

	_, err := uc.Login(models.User{
		Email:    "john@example.com",
		Password: "password123",
	})

	assert.EqualError(t, err, "incorrect password")
}

func TestPromote_Success(t *testing.T) {
	repo, _, _, _, _, uc := setup()

	repo.On("FindByEmail", "john@example.com").Return(models.User{}, nil)
	repo.On("UpdateRole", "john@example.com", "admin").Return(nil)

	err := uc.Promote("john@example.com")
	assert.NoError(t, err)
}

func TestPromote_super_admin(t *testing.T) {
	repo, _, _, _, _, uc := setup()
	user := models.User{
		Email: "john@example.com",
		Role:  "superadmin",
	}

	repo.On("FindByEmail", "john@example.com").Return(user, nil)

	err := uc.Promote("john@example.com")
	assert.EqualError(t, err, "superadmin cannot be promoted")
}

func TestDemote_UserNotFound(t *testing.T) {
	repo, _, _, _, _, uc := setup()

	repo.On("FindByEmail", "john@example.com").Return(models.User{}, errors.New("not found"))

	err := uc.Demote("john@example.com")
	assert.EqualError(t, err, "user not found")
}

func TestDemote_success(t *testing.T) {
	repo, _, _, _, _, uc := setup()
	user := models.User{
		Email: "john@example.com",
		Role:  "admin",
	}

	repo.On("FindByEmail", "john@example.com").Return(user, nil)
	repo.On("UpdateRole", "john@example.com", "user").Return(nil)

	err := uc.Demote("john@example.com")
	assert.NoError(t, err)
}

func TestDemote_super_admin(t *testing.T) {
	repo, _, _, _, _, uc := setup()
	user := models.User{
		Email: "john@example.com",
		Role:  "superadmin",
	}

	repo.On("FindByEmail", "john@example.com").Return(user, nil)

	err := uc.Demote("john@example.com")
	assert.EqualError(t, err, "superadmin cannot be demoted")
}

func TestRefreshToken_ExpiredToken(t *testing.T) {
	_, _, mockTokenService, mockTokenRepo, _, uc := setup()

	refreshStr := "refresh.jwt.token"
	expiredTime := time.Now().Add(-time.Hour)

	mockTokenService.On("VerifyRefreshToken", refreshStr).Return(
		&models.UserRefreshClaims{
			UserID:    "123",
			Email:     "test@example.com",
			Role:      "user",
			TokenID:   "token123",
			ExpiresAt: expiredTime,
		}, nil,
	)

	mockTokenRepo.On("DeleteToken", "token123").Return(nil)

	_, err := uc.RefreshToken(refreshStr)

	assert.Error(t, err)
	assert.EqualError(t, err, "the refresh token expired")

	mockTokenService.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	_, mockHasher, mockTokenService, mockTokenRepo, _, uc := setup()

	refreshStr := "refresh.jwt.token"
	validTime := time.Now().Add(time.Hour)

	mockTokenService.On("VerifyRefreshToken", refreshStr).Return(
		&models.UserRefreshClaims{
			UserID:    "123",
			Email:     "test@example.com",
			Role:      "user",
			TokenID:   "token123",
			ExpiresAt: validTime,
		}, nil,
	)

	mockTokenRepo.On("GetToken", "token123").Return(
		&models.Token{
			ID:    "token123",
			Token: "hashed_token",
		}, nil,
	)

	mockTokenService.On("VerifyToken", "hashed_token", refreshStr).Return(true)

	mockTokenService.On("GenerateAccessToken", "123", "test@example.com", "user").Return("new_access_token", nil)

	token, err := uc.RefreshToken(refreshStr)

	assert.NoError(t, err)
	assert.Equal(t, "new_access_token", token)

	mockTokenService.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}
func TestLogout_Success(t *testing.T) {
	_, _, mockTokenService, mockTokenRepo, _, uc := setup()

	refreshStr := "refresh.jwt.token"
	validTime := time.Now().Add(time.Hour)

	// Step 1: Mock VerifyRefreshToken to return valid claims
	mockTokenService.On("VerifyRefreshToken", refreshStr).Return(
		&models.UserRefreshClaims{
			UserID:    "123",
			Email:     "test@example.com",
			Role:      "user",
			TokenID:   "token123",
			ExpiresAt: validTime,
		}, nil,
	)

	// Step 2: Mock GetToken to return stored hashed token
	mockTokenRepo.On("GetToken", "token123").Return(
		&models.Token{
			ID:    "token123",
			Token: "hashed_token",
		}, nil,
	)

	// Step 3: Mock VerifyToken to match
	mockTokenService.On("VerifyToken", "hashed_token", refreshStr).Return(true)

	// Step 4: Mock DeleteToken success
	mockTokenRepo.On("DeleteToken", "token123").Return(nil)

	err := uc.Logout(refreshStr)

	assert.NoError(t, err)
	mockTokenService.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestLogout_VerifyRefreshTokenError(t *testing.T) {
	// Arrange
	_, _, mockTokenService, _, _, uc := setup()

	refreshStr := "invalid.token"

	// Return nil pointer of correct type to avoid panic
	mockTokenService.
		On("VerifyRefreshToken", refreshStr).
		Return((*models.UserRefreshClaims)(nil), errors.New("invalid token"))

	// Act
	err := uc.Logout(refreshStr)

	// Assert
	assert.EqualError(t, err, "invalid token")
	mockTokenService.AssertExpectations(t)
}

func TestLogout_ExpiredToken(t *testing.T) {
	_, _, mockTokenService, mockTokenRepo, _, uc := setup()

	refreshStr := "refresh.jwt.token"
	expiredTime := time.Now().Add(-time.Hour)

	mockTokenService.On("VerifyRefreshToken", refreshStr).Return(
		&models.UserRefreshClaims{
			TokenID:   "token123",
			ExpiresAt: expiredTime,
		}, nil,
	)

	mockTokenRepo.On("DeleteToken", "token123").Return(nil)

	err := uc.Logout(refreshStr)
	assert.EqualError(t, err, "the refresh token expired")
	mockTokenService.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestLogout_TokenNotFoundInDB(t *testing.T) {
	_, _, mockTokenService, mockTokenRepo, _, uc := setup()

	refreshStr := "refresh.jwt.token"
	validTime := time.Now().Add(time.Hour)

	mockTokenService.On("VerifyRefreshToken", refreshStr).Return(
		&models.UserRefreshClaims{
			TokenID:   "token123",
			ExpiresAt: validTime,
		}, nil,
	)

	mockTokenRepo.On("GetToken", "token123").Return(&models.Token{}, errors.New("token not found"))

	err := uc.Logout(refreshStr)
	assert.EqualError(t, err, "token not found")
	mockTokenService.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestLogout_InvalidTokenHash(t *testing.T) {
	_, _, mockTokenService, mockTokenRepo, _, uc := setup()

	refreshStr := "refresh.jwt.token"
	validTime := time.Now().Add(time.Hour)

	mockTokenService.On("VerifyRefreshToken", refreshStr).Return(
		&models.UserRefreshClaims{
			TokenID:   "token123",
			ExpiresAt: validTime,
		}, nil,
	)

	mockTokenRepo.On("GetToken", "token123").Return(
		&models.Token{
			ID:    "token123",
			Token: "hashed_token",
		}, nil,
	)

	mockTokenService.On("VerifyToken", "hashed_token", refreshStr).Return(false)

	err := uc.Logout(refreshStr)
	assert.EqualError(t, err, "invalid refresh token")
	mockTokenService.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}
