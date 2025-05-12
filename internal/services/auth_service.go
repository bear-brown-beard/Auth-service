package services

import (
	"context"
	"time"

	"auth-service/internal/models"
	"auth-service/internal/repositories"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the interface for authentication service operations
type AuthService interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	userRepository repositories.AuthRepository
	jwtSecret      []byte
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepository repositories.AuthRepository, jwtSecret string) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtSecret:      []byte(jwtSecret),
	}
}

// Register creates a new user with hashed password
func (s *authService) Register(ctx context.Context, user *models.User) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	// Create user in database
	return s.userRepository.CreateUser(ctx, user)
}

// Login verifies user credentials and returns JWT token
func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// Get user from database
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", models.ErrUserNotFound
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", models.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Sign token
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
