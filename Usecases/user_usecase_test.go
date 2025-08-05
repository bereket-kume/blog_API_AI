package usecases

import (
	"blog-api/Domain/models"
	"blog-api/mocks" // generated mocks
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (*mocks.UserRepository, *mocks.MockHasher, *mocks.MockTokenService, *mocks.MockTokenRepository, UserUsecaseInterface) {
	repo := new(mocks.UserRepository)
	hasher := new(mocks.MockHasher)
	tokenSvc := new(mocks.MockTokenService)
	tokenRepo := new(mocks.MockTokenRepository)

	uc := NewUserUsecase(repo, hasher, tokenSvc, tokenRepo)
	return repo, hasher, tokenSvc, tokenRepo, uc
}

func TestRegister_Success_AdminFirstUser(t *testing.T) {
	repo, hasher, _, _, uc := setup()

	user := models.User{
		Username: "john",
		Email:    "john@example.com",
		Password: "password123",
	}

	hasher.On("HashPassword", "password123").Return("hashed_pw")
	repo.On("CountUsers").Return(0, nil) // first user
	repo.On("Insert", mock.Anything).Return(nil)

	err := uc.Register(user)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}

func TestRegister_InvalidEmail(t *testing.T) {
	_, _, _, _, uc := setup()

	user := models.User{
		Username: "john",
		Email:    "invalid_email",
		Password: "password123",
	}

	err := uc.Register(user)
	assert.EqualError(t, err, "Invalid Email address")
}

func TestLogin_Success(t *testing.T) {
	repo, hasher, tokenSvc, tokenRepo, uc := setup()

	existing := models.User{
		ID:       "user1",
		Email:    "john@example.com",
		Password: "hashed_pw",
		Role:     "user",
		Verified: true,
	}

	repo.On("FindByEmail", "john@example.com").Return(existing, nil)
	hasher.On("VerifyPassword", "hashed_pw", "password123").Return(true)
	tokenSvc.On("GenerateAccessToken", existing.ID, existing.Email, existing.Role).Return("access123", nil)
	tokenSvc.On("GenerateRefreshToken", existing.ID, existing.Email, existing.Role).
		Return(&models.Token{Token: "refresh_plain"}, nil)
	hasher.On("HashPassword", "refresh_plain").Return("refresh_hashed")
	tokenRepo.On("CreateToken", mock.Anything).Return(nil)

	out, err := uc.Login(models.User{
		Email:    "john@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.Equal(t, "access123", out.Access_token)
	assert.Equal(t, "refresh_hashed", out.Refresh_token)
}

func TestLogin_UserNotVerified(t *testing.T) {
	repo, _, _, _, uc := setup()

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

func TestPromote_Success(t *testing.T) {
	repo, _, _, _, uc := setup()

	repo.On("FindByEmail", "john@example.com").Return(models.User{}, nil)
	repo.On("UpdateRole", "john@example.com", "admin").Return(nil)

	err := uc.Promote("john@example.com")
	assert.NoError(t, err)
}

func TestPromote_super_admin(t *testing.T) {
	repo, _, _, _, uc := setup()
	user := models.User{
		Email:    "john@example.com",
		Password: "password123",
		Role:     "superadmin",
	}

	repo.On("FindByEmail", "john@example.com").Return(user, nil)

	err := uc.Promote("john@example.com")
	assert.EqualError(t, err, "superadmin cannot be promoted")
}

func TestDemote_UserNotFound(t *testing.T) {
	repo, _, _, _, uc := setup()

	repo.On("FindByEmail", "john@example.com").Return(models.User{}, errors.New("not found"))

	err := uc.Demote("john@example.com")
	assert.EqualError(t, err, "user not found")
}
func TestDemote_success(t *testing.T) {
	repo, _, _, _, uc := setup()
	user := models.User{
		Email:    "john@example.com",
		Password: "password123",
		Role:     "admin",
	}

	repo.On("FindByEmail", "john@example.com").Return(user, nil)
	repo.On("UpdateRole", "john@example.com", "user").Return(nil)

	err := uc.Demote("john@example.com")
	assert.NoError(t, err)
}

func TestDemote_super_admin(t *testing.T) {
	repo, _, _, _, uc := setup()
	user := models.User{
		Email:    "john@example.com",
		Password: "password123",
		Role:     "superadmin",
	}

	repo.On("FindByEmail", "john@example.com").Return(user, nil)

	err := uc.Demote("john@example.com")
	assert.EqualError(t, err, "superadmin cannot be demoted")
}
