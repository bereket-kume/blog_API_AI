package services

import (
	"blog-api/Domain/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	accessSecretKey  string
	refreshSecretKey string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

func NewJWTService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWTService {
	return &JWTService{
		accessSecretKey:  accessSecret,
		refreshSecretKey: refreshSecret,
		accessTokenTTL:   accessTTL,
		refreshTokenTTL:  refreshTTL,
	}
}

func (j *JWTService) GenerateAccessToken(userID, email, role string) (string, error) {
	exp := time.Now().Add(j.accessTokenTTL).Unix()

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.accessSecretKey))

	return signed, err
}

func (j *JWTService) GenerateRefreshToken(userID, email, role string) (*models.Token, error) {
	exp := time.Now().Add(j.refreshTokenTTL)
	iat := time.Now()

	tokenID := uuid.New().String()

	claims := jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"role":     role,
		"token_id": tokenID,
		"exp":      exp.Unix(),
		"iat":      iat.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.refreshSecretKey))
	if err != nil {
		return nil, err
	}

	return &models.Token{
		ID:        tokenID,
		Token:     signed,
		ExpiresAt: exp,
		CreatedAt: iat,
		UserID:    userID,
		// IP: '',
		// Device: '',
	}, nil
}
func (j *JWTService) VerifyAccessToken(tokenStr string) (*models.UserAccessClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.accessSecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	return &models.UserAccessClaims{
		UserID: claims["user_id"].(string),
		Email:  claims["email"].(string),
		Role:   claims["role"].(string),
	}, nil
}

func (j *JWTService) VerifyRefreshToken(tokenStr string) (*models.UserRefreshClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.refreshSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	claims := token.Claims.(jwt.MapClaims)
	expUnix := int64(claims["exp"].(float64))
	iatUnix := int64(claims["iat"].(float64))
	return &models.UserRefreshClaims{
		UserID:    claims["user_id"].(string),
		Email:     claims["email"].(string),
		Role:      claims["role"].(string),
		TokenID:   claims["token_id"].(string),
		ExpiresAt: time.Unix(expUnix, 0),
		CreatedAt: time.Unix(iatUnix, 0),
	}, nil

}

func (j *JWTService) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (j *JWTService) VerifyToken(hashed, token string) bool {
	return hashed == j.HashToken(token)
}
