package mocks

import (
	"blog-api/Domain/models"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateAccessToken(userID, email, role string) (string, error) {
	args := m.Called(userID, email, role)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken(userID, email, role string) (*models.Token, error) {
	args := m.Called(userID, email, role)
	return args.Get(0).(*models.Token), args.Error(1)
}

func (m *MockTokenService) VerifyAccessToken(tokenStr string) (*models.UserAccessClaims, error) {
	args := m.Called(tokenStr)
	return args.Get(0).(*models.UserAccessClaims), args.Error(1)
}

func (m *MockTokenService) VerifyRefreshToken(tokenStr string) (*models.UserRefreshClaims, error) {
	args := m.Called(tokenStr)
	return args.Get(0).(*models.UserRefreshClaims), args.Error(1)
}

func (m *MockTokenService) GenerateRandomJWT(expiredAt time.Duration) (*models.Token, error) {
	args := m.Called(expiredAt)
	return args.Get(0).(*models.Token), args.Error(1)
}

func (m *MockTokenService) VerifyJWT(tokenStr string) (models.TokenClaims, error) {
	args := m.Called(tokenStr)
	return args.Get(0).(models.TokenClaims), args.Error(1)
}

func (m *MockTokenService) HashToken(token string) string {
	args := m.Called(token)
	return args.String(0)
}

func (m *MockTokenService) VerifyToken(hashed, token string) bool {
	args := m.Called(hashed, token)
	return args.Bool(0)
}
