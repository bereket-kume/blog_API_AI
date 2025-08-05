package interfaces

import (
	"blog-api/Domain/models"
)

type TokenService interface {
	GenerateAccessToken(userID, email, role string) (string, error)
	GenerateRefreshToken(userID, email, role string) (*models.Token, error)
	VerifyAccessToken(tokenStr string) (*models.UserAccessClaims, error)
	VerifyRefreshToken(tokenStr string) (*models.UserRefreshClaims, error)
}
type TokenRepository interface {
	CreateToken(token *models.Token) error
	DeleteToken(tokenID string) error
	Update(token *models.Token) error
	GetToken(tokenID string) (*models.Token, error)
}
