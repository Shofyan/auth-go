package usecase

import (
	"context"
	"time"

	"auth-go/internal/application/dto"
	"auth-go/internal/domain/entity"
	"auth-go/internal/domain/repository"
	"auth-go/internal/domain/service"
	apperrors "auth-go/pkg/errors"
)

// RefreshTokenUseCase handles token refresh with rotation
type RefreshTokenUseCase struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	tokenService     service.TokenService
}

// NewRefreshTokenUseCase creates a new refresh token use case
func NewRefreshTokenUseCase(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	tokenService service.TokenService,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		tokenService:     tokenService,
	}
}

// Execute executes the refresh token use case with token rotation
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Find refresh token
	refreshToken, err := uc.refreshTokenRepo.FindByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	// Check if token is revoked (potential token reuse attack)
	if refreshToken.IsRevoked {
		// Revoke all tokens in this family as a security measure
		_ = uc.refreshTokenRepo.RevokeByTokenFamily(ctx, refreshToken.TokenFamily)
		return nil, apperrors.ErrTokenReuse
	}

	// Check if token is expired
	if refreshToken.IsExpired() {
		return nil, apperrors.ErrExpiredToken
	}

	// Revoke current token (token rotation)
	refreshToken.Revoke()
	if err := uc.refreshTokenRepo.Update(ctx, refreshToken); err != nil {
		return nil, err
	}

	// Get user
	user, err := uc.userRepo.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperrors.ErrUserInactive
	}

	// Generate new access token
	accessToken, err := uc.tokenService.GenerateAccessToken(service.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  user.Roles,
	})
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshTokenStr, err := uc.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Create new refresh token entity (same token family, different token)
	expiresAt := time.Now().Add(uc.tokenService.GetRefreshTokenExpiry())
	newRefreshToken := entity.NewRefreshToken(user.ID, newRefreshTokenStr, expiresAt, refreshToken.TokenFamily)
	newRefreshToken.ParentToken = &refreshToken.Token // Track parent for rotation chain

	// Save new refresh token
	if err := uc.refreshTokenRepo.Create(ctx, newRefreshToken); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    int64(uc.tokenService.GetAccessTokenExpiry().Seconds()),
	}, nil
}
