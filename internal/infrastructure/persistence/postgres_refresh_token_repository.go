package persistence

import (
	"context"
	"database/sql"
	"errors"

	"auth-go/internal/domain/entity"
	"auth-go/internal/domain/repository"
	apperrors "auth-go/pkg/errors"

	"github.com/google/uuid"
)

// PostgresRefreshTokenRepository implements RefreshTokenRepository using PostgreSQL
type PostgresRefreshTokenRepository struct {
	db *sql.DB
}

// NewPostgresRefreshTokenRepository creates a new PostgreSQL refresh token repository
func NewPostgresRefreshTokenRepository(db *sql.DB) repository.RefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *PostgresRefreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, is_revoked, token_family, parent_token)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
		token.IsRevoked,
		token.TokenFamily,
		token.ParentToken,
	)

	return err
}

// FindByToken finds a refresh token by token string
func (r *PostgresRefreshTokenRepository) FindByToken(ctx context.Context, tokenStr string) (*entity.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, is_revoked, revoked_at, token_family, parent_token
		FROM refresh_tokens
		WHERE token = $1
	`

	token := &entity.RefreshToken{}
	var revokedAt sql.NullTime
	var parentToken sql.NullString

	err := r.db.QueryRowContext(ctx, query, tokenStr).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.IsRevoked,
		&revokedAt,
		&token.TokenFamily,
		&parentToken,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrInvalidToken
		}
		return nil, err
	}

	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}

	if parentToken.Valid {
		token.ParentToken = &parentToken.String
	}

	return token, nil
}

// FindByUserID finds all refresh tokens for a user
func (r *PostgresRefreshTokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, is_revoked, revoked_at, token_family, parent_token
		FROM refresh_tokens
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entity.RefreshToken
	for rows.Next() {
		token := &entity.RefreshToken{}
		var revokedAt sql.NullTime
		var parentToken sql.NullString

		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.IsRevoked,
			&revokedAt,
			&token.TokenFamily,
			&parentToken,
		)
		if err != nil {
			return nil, err
		}

		if revokedAt.Valid {
			token.RevokedAt = &revokedAt.Time
		}

		if parentToken.Valid {
			token.ParentToken = &parentToken.String
		}

		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

// Update updates a refresh token
func (r *PostgresRefreshTokenRepository) Update(ctx context.Context, token *entity.RefreshToken) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = $2, revoked_at = $3
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, token.ID, token.IsRevoked, token.RevokedAt)
	return err
}

// RevokeByTokenFamily revokes all tokens in a token family
func (r *PostgresRefreshTokenRepository) RevokeByTokenFamily(ctx context.Context, tokenFamily uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true, revoked_at = NOW()
		WHERE token_family = $1 AND is_revoked = false
	`

	_, err := r.db.ExecContext(ctx, query, tokenFamily)
	return err
}

// RevokeByUserID revokes all tokens for a user
func (r *PostgresRefreshTokenRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true, revoked_at = NOW()
		WHERE user_id = $1 AND is_revoked = false
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// DeleteExpired deletes all expired tokens
func (r *PostgresRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	_, err := r.db.ExecContext(ctx, query)
	return err
}
