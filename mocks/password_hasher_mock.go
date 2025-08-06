package mocks

import "github.com/stretchr/testify/mock"

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), nil
}

func (m *MockHasher) VerifyPassword(hashed, password string) bool {
	args := m.Called(hashed, password)
	return args.Bool(0)
}
