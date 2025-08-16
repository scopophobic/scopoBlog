package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService struct {
	JWTSecret         string
	AdminPasswordHash string
}

func NewAuthService(jwtSecret, adminPasswordHash string) *AuthService {
	return &AuthService{
		JWTSecret:         jwtSecret,
		AdminPasswordHash: adminPasswordHash,
	}
}

// GeneratePasswordHash creates a bcrypt hash from a plain text password
func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword checks if a password matches the stored hash
func (a *AuthService) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.AdminPasswordHash), []byte(password))
	return err == nil
}

// GenerateJWT creates a new JWT token for the admin user
func (a *AuthService) GenerateJWT() (string, error) {
	claims := Claims{
		Username: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.JWTSecret))
}

// ValidateJWT validates a JWT token and returns the claims
func (a *AuthService) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
