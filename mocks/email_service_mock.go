// mocks/email_mock.go
package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func (m *MockEmailService) SendVerificationEmail(username, email, token string) error {
	args := m.Called(username, email, token)
	return args.Error(0)
}

func (m *MockEmailService) SendPasswordResetEmail(username, email, token string) error {
	args := m.Called(username, email, token)
	return args.Error(0)
}
