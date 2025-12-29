package usecase

import (
	"context"

	"auth-go/internal/domain/repository"

	"github.com/google/uuid"
)

// LogoutUseCase handles user logout by revoking refresh tokens
type LogoutUseCase struct {
	refreshTokenRepo repository.RefreshTokenRepository
}

// NewLogoutUseCase creates a new logout use case
func NewLogoutUseCase(refreshTokenRepo repository.RefreshTokenRepository) *LogoutUseCase {
	return &LogoutUseCase{
		refreshTokenRepo: refreshTokenRepo,
	}
}

// Execute executes the logout use case (revokes all user tokens)
func (uc *LogoutUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	return uc.refreshTokenRepo.RevokeByUserID(ctx, userID)
}
