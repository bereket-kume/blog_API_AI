package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptHasher(t *testing.T) {
	hasher := BcryptHasher{}

	// 1️⃣ Test hashing
	password := "mysecretpassword"
	hash, err := hasher.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash) // hash should be different from password

	// 2️⃣ Test correct password verification
	isValid := hasher.VerifyPassword(hash, password)
	assert.True(t, isValid, "expected password to be valid")

	// 3️⃣ Test incorrect password verification
	isInvalid := hasher.VerifyPassword(hash, "wrongpassword")
	assert.False(t, isInvalid, "expected password to be invalid")
}
