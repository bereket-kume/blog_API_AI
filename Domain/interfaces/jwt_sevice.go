package interfaces

import (
	"blog-api/Domain/models"
)

type TokenService interface {
	GenerateToken(userID, email, role string) (*models.Token, error)
	VerifyToken(tokenStr string) (*models.UserClaims, error)
}
type TokenRepository interface {
	CreateToken(token models.Token) error
	DeleteToken(tokenID string) error
	Update(token models.Token) error
	GetToken(tokenID string) (models.Token, error)
}
