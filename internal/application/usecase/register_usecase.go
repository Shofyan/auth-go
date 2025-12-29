package usecase

import (
	"context"

	"auth-go/internal/application/dto"
	"auth-go/internal/domain/entity"
	"auth-go/internal/domain/repository"
	"auth-go/internal/domain/service"
	"auth-go/internal/domain/valueobject"
	apperrors "auth-go/pkg/errors"
)

// RegisterUseCase handles user registration
type RegisterUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher service.PasswordHasher
}

// NewRegisterUseCase creates a new register use case
func NewRegisterUseCase(
	userRepo repository.UserRepository,
	passwordHasher service.PasswordHasher,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

// Execute executes the register use case
func (uc *RegisterUseCase) Execute(ctx context.Context, req dto.RegisterRequest) error {
	// Validate email
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return apperrors.ErrInvalidEmail
	}

	// Check if user already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, email.Value())
	if err != nil {
		return err
	}
	if exists {
		return apperrors.ErrUserAlreadyExists
	}

	// Validate password
	password, err := valueobject.NewPassword(req.Password)
	if err != nil {
		return apperrors.ErrInvalidPassword
	}

	// Hash password
	passwordHash, err := uc.passwordHasher.Hash(password.Value())
	if err != nil {
		return err
	}

	// Create user entity
	user := entity.NewUser(email.Value(), passwordHash)

	// Save user
	return uc.userRepo.Create(ctx, user)
}
