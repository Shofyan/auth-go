package service

import (
	"time"

	"auth-go/internal/domain/entity"

	"github.com/google/uuid"
)

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID uuid.UUID
	Email  string
	Roles  []entity.Role
}

// TokenPair represents an access and refresh token pair
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds
}

// TokenService defines the interface for token operations
type TokenService interface {
	// GenerateAccessToken generates a JWT access token
	GenerateAccessToken(claims TokenClaims) (string, error)

	// GenerateRefreshToken generates a refresh token
	GenerateRefreshToken() (string, error)

	// ValidateAccessToken validates and parses an access token
	ValidateAccessToken(token string) (*TokenClaims, error)

	// GetAccessTokenExpiry returns the access token expiry duration
	GetAccessTokenExpiry() time.Duration

	// GetRefreshTokenExpiry returns the refresh token expiry duration
	GetRefreshTokenExpiry() time.Duration
}
