package mocks

import (
	"blog-api/Domain/models"

	"github.com/stretchr/testify/mock"
)

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) CreateToken(token *models.Token) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteToken(tokenID string) error {
	args := m.Called(tokenID)
	return args.Error(0)
}

func (m *MockTokenRepository) Update(token *models.Token) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenRepository) GetToken(tokenID string) (*models.Token, error) {
	args := m.Called(tokenID)
	return args.Get(0).(*models.Token), args.Error(1)
}
