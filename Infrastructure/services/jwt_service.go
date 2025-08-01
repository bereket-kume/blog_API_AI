package services

import (
	"blog-api/Domain/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = []byte("asfm")

type JWTService struct {
	SecretKey string
}

func NewJWTService() *JWTService {
	return &JWTService{SecretKey: string(SecretKey)}
}

func (j *JWTService) GenerateToken(userID, email, role string) (models.Token, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTService) VerifyToken(tokenStr string) (*domain.UserClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, domain.ErrTokenInvalid
	}
	claims := token.Claims.(jwt.MapClaims)
	return &domain.UserClaims{
		UserID: claims["user_id"].(string),
		Email:  claims["email"].(string),
		Role:   claims["role"].(string),
	}, nil
}
