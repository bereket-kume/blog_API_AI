package db_models

import (
	"blog-api/Domain/models"
	"time"
)

type Token struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	TokenHash string    `bson:"token_hash"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	IP        string    `bson:"ip"`
	Device    string    `bson:"device"`
	Email     string    `bson:"email"`
}

func FromDomainToken(token *models.Token) *Token {
	return &Token{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.Token,
		CreatedAt: token.CreatedAt,
		ExpiresAt: token.ExpiresAt,
		IP:        token.IP,
		Device:    token.Device,
		Email:     token.Email,
	}
}

func ToDomainToken(token *Token) *models.Token {
	return &models.Token{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.TokenHash,
		CreatedAt: token.CreatedAt,
		ExpiresAt: token.ExpiresAt,
		IP:        token.IP,
		Device:    token.Device,
		Email:     token.Email,
	}
}
