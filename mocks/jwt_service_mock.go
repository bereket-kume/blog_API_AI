package mocks

import (
	"blog-api/Domain/models"

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
