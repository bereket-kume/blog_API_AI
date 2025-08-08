package interfaces

import (
	"blog-api/Domain/models"
	"time"
)

type TokenService interface {
	GenerateAccessToken(userID, email, role string) (string, error)
	GenerateRefreshToken(userID, email, role string) (*models.Token, error)
	VerifyAccessToken(tokenStr string) (*models.UserAccessClaims, error)
	VerifyRefreshToken(tokenStr string) (*models.UserRefreshClaims, error)
	GenerateRandomJWT(expiredAt time.Duration) (*models.Token, error)
	VerifyJWT(tokenStr string) (models.TokenClaims, error)
	HashToken(token string) string
	VerifyToken(hashed, token string) bool
}
type TokenRepository interface {
	CreateToken(token *models.Token) error
	DeleteToken(tokenID string) error
	Update(token *models.Token) error
	GetToken(tokenID string) (*models.Token, error)
}
