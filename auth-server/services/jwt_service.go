package services

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) 

type JWTService interface {
	GenerateToken(userID string) (string, string)
	ValidateToken(tokenString string) (map[string]interface{}, string)
}

type jwtService struct {
	jwtSecret []byte
}

func NewJWTService() JWTService {
	return &jwtService{
		jwtSecret: jwtSecret,
	}
}

func (s *jwtService) GenerateToken(userID string) (string, string) {
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id":    userID,
		"created_at": now.Unix(),
		"exp":        now.Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "TOKEN_GENERATION_ERROR"
	}

	return tokenString, ""
}

func (s *jwtService) ValidateToken(tokenString string) (map[string]interface{}, string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, "INVALID_TOKEN"
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, "INVALID_TOKEN"
	}

	return claims, ""
}
