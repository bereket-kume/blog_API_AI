package models

type Hasher interface {
	HashPassword(password string) string
	VerifyPassword(hashed, password string) bool
}

type UserClaims struct {
	UserId string
	Email  string
	Role   string
}

type TokenService interface {
	GenerateToken(userID, email, role string) (string, error)
	VerifyToken(tokenStr string) (*UserClaims, error)
}
