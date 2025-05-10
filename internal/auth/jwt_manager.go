package auth

import (
	"time"
	"github.com/dgrijalva/jwt-go"
)

// JWTManager handles JWT token generation and verification
type JWTManager struct {
	secretKey string
	issuer    string
	expiry    time.Duration
}

// NewJWTManager creates a new JWTManager instance
func NewJWTManager(secretKey, issuer string, expiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		issuer:    issuer,
		expiry:    expiry,
	}
}

// GenerateToken generates a new JWT token
func (m *JWTManager) GenerateToken(email string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   email,
		ExpiresAt: time.Now().Add(m.expiry).Unix(),
		Issuer:    m.issuer,
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// ParseToken verifies and parses a JWT token
func (m *JWTManager) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})
}

// GetDefaultJWTManager returns a JWTManager with default settings
func GetDefaultJWTManager(secretKey string) *JWTManager {
	return NewJWTManager(
		secretKey,
		"auth-service",
		24*time.Hour, // Default token expiry: 24 hours
	)
}
