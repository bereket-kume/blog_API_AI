// Mocks/mock_user_repository.go
package mocks

import (
	"blog-api/Domain/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Insert(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePass(email string, passwordHash string) error {
	args := m.Called(email, passwordHash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateRole(email string, role string) error {
	args := m.Called(email, role)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserRepository) Verify(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserRepository) CountUsers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
