package repositories

import (
	"auth-service/internal/models"
	"context"
	"database/sql"
)

// AuthRepository defines the interface for authentication-related database operations
type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type authRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new instance of AuthRepository
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

// CreateUser inserts a new user into the database
func (r *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	return r.db.QueryRowContext(ctx, query,
		user.Email,
		user.PasswordHash,
	).Scan(&id)
}

// GetUserByEmail retrieves a user by their email address
func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`

	var user models.User
	row := r.db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Return nil user instead of error for non-existent users
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
