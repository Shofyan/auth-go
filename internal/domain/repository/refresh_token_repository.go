package repository

import (
	"context"

	"auth-go/internal/domain/entity"

	"github.com/google/uuid"
)

// RefreshTokenRepository defines the interface for refresh token persistence
type RefreshTokenRepository interface {
	// Create creates a new refresh token
	Create(ctx context.Context, token *entity.RefreshToken) error

	// FindByToken finds a refresh token by token string
	FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error)

	// FindByUserID finds all refresh tokens for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.RefreshToken, error)

	// Update updates a refresh token
	Update(ctx context.Context, token *entity.RefreshToken) error

	// RevokeByTokenFamily revokes all tokens in a token family (for rotation security)
	RevokeByTokenFamily(ctx context.Context, tokenFamily uuid.UUID) error

	// RevokeByUserID revokes all tokens for a user
	RevokeByUserID(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired deletes all expired tokens
	DeleteExpired(ctx context.Context) error
}
