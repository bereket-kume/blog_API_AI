package models

import (
	"time"
)

type Token struct {
	ID        string // Use string instead of ObjectID
	UserID    string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	IP        string
	Device    string
}

type UserAccessClaims struct {
	UserID    string
	Email     string
	Role      string
	ExpiresAt time.Time
	CreatedAt time.Time
}
type UserRefreshClaims struct {
	UserID    string
	Email     string
	Role      string
	TokenID   string
	ExpiresAt time.Time
	CreatedAt time.Time
}
