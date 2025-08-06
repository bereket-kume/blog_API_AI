package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupJWTService() *JWTService {
	return NewJWTService(
		"test-access-secret",
		"test-refresh-secret",
		15*time.Minute,
		7*24*time.Hour,
	)
}

func TestGenerateAccessToken(t *testing.T) {
	svc := setupJWTService()

	tokenStr, err := svc.GenerateAccessToken("user123", "test@example.com", "admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := svc.VerifyAccessToken(tokenStr)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "admin", claims.Role)
}

func TestGenerateRefreshToken(t *testing.T) {
	svc := setupJWTService()

	token, err := svc.GenerateRefreshToken("user123", "test@example.com", "admin")
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.ID)
	assert.NotEmpty(t, token.Token)
	assert.WithinDuration(t, time.Now().Add(7*24*time.Hour), token.ExpiresAt, time.Minute)
	assert.Equal(t, "user123", token.UserID)

	claims, err := svc.VerifyRefreshToken(token.Token)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, token.ID, claims.TokenID)
}

func TestVerifyAccessToken_Invalid(t *testing.T) {
	svc := setupJWTService()

	_, err := svc.VerifyAccessToken("invalid.token.value")
	assert.Error(t, err)
}

func TestVerifyRefreshToken_Invalid(t *testing.T) {
	svc := setupJWTService()

	_, err := svc.VerifyRefreshToken("invalid.token.value")
	assert.Error(t, err)
}
