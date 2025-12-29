package persistence

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"auth-go/internal/domain/entity"
	"auth-go/internal/domain/repository"
	apperrors "auth-go/pkg/errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, roles, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.String()
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		pq.Array(roles),
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return apperrors.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

// FindByID finds a user by ID
func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, roles, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`

	user := &entity.User{}
	var roles pq.StringArray
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&roles,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	user.Roles = make([]entity.Role, len(roles))
	for i, role := range roles {
		user.Roles[i] = entity.ParseRole(role)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByEmail finds a user by email
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, roles, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`

	user := &entity.User{}
	var roles pq.StringArray
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&roles,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	user.Roles = make([]entity.Role, len(roles))
	for i, role := range roles {
		user.Roles[i] = entity.ParseRole(role)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// Update updates a user
func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, roles = $4, is_active = $5, updated_at = $6, last_login_at = $7
		WHERE id = $1
	`

	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.String()
	}

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		pq.Array(roles),
		user.IsActive,
		user.UpdatedAt,
		user.LastLoginAt,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

// ExistsByEmail checks if a user exists by email
func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}
