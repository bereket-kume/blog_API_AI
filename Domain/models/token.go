package models

import "time"

type Token struct {
	ID        string // Use string instead of ObjectID
	UserID    string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	IP        string
	Device    string
}

type UserClaims struct {
	UserId string
	Email  string
	Role   string
}
