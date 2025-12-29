package entity

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token for token rotation
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	IsRevoked bool
	RevokedAt *time.Time
	// Token family for rotation detection
	TokenFamily uuid.UUID
	// Previous token in the rotation chain (for detecting reuse)
	ParentToken *string
}

// NewRefreshToken creates a new refresh token
func NewRefreshToken(userID uuid.UUID, token string, expiresAt time.Time, tokenFamily uuid.UUID) *RefreshToken {
	return &RefreshToken{
		ID:          uuid.New(),
		UserID:      userID,
		Token:       token,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
		IsRevoked:   false,
		TokenFamily: tokenFamily,
	}
}

// IsExpired checks if the token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked
}

// Revoke revokes the refresh token
func (rt *RefreshToken) Revoke() {
	rt.IsRevoked = true
	now := time.Now()
	rt.RevokedAt = &now
}
