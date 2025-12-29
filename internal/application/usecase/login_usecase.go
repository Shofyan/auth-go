package usecase

import (
	"context"
	"log"
	"time"

	"auth-go/internal/application/dto"
	"auth-go/internal/domain/entity"
	"auth-go/internal/domain/repository"
	"auth-go/internal/domain/service"
	apperrors "auth-go/pkg/errors"

	"github.com/google/uuid"
)

// LoginUseCase handles user login with JWT access and refresh tokens
type LoginUseCase struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	passwordHasher   service.PasswordHasher
	tokenService     service.TokenService
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	passwordHasher service.PasswordHasher,
	tokenService service.TokenService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		passwordHasher:   passwordHasher,
		tokenService:     tokenService,
	}
}

// Execute executes the login use case
func (uc *LoginUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperrors.ErrUserInactive
	}

	// Verify password
	if err := uc.passwordHasher.Compare(req.Password, user.PasswordHash); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Update last login
	user.UpdateLastLogin()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		// Log error but don't fail the login
		log.Printf("Failed to update last login for user %s: %v", user.ID, err)
	}

	// Generate access token
	accessToken, err := uc.tokenService.GenerateAccessToken(service.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  user.Roles,
	})
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenStr, err := uc.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Create refresh token entity with new token family
	tokenFamily := uuid.New()
	expiresAt := time.Now().Add(uc.tokenService.GetRefreshTokenExpiry())
	refreshToken := entity.NewRefreshToken(user.ID, refreshTokenStr, expiresAt, tokenFamily)

	// Save refresh token
	if err := uc.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    int64(uc.tokenService.GetAccessTokenExpiry().Seconds()),
	}, nil
}
