package mocks

import (
	"blog-api/Domain/models"

	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Insert(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) FindByEmail(email string) (models.User, error) {
	args := m.Called(email)
	user, _ := args.Get(0).(models.User)
	return user, args.Error(1)
}

func (m *UserRepository) UpdatePass(email, passowrdHash string) error {
	args := m.Called(email, passowrdHash)
	return args.Error(0)
}

func (m *UserRepository) UpdateRole(email, role string) error {
	args := m.Called(email, role)
	return args.Error(0)
}

func (m *UserRepository) Delete(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *UserRepository) Verify(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *UserRepository) CountUsers() (int64, error) {
	args := m.Called()
	count, _ := args.Get(0).(int64)
	return count, args.Error(1)
}
