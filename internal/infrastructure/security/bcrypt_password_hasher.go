package security

import (
	"auth-go/internal/domain/service"

	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher implements PasswordHasher using bcrypt
type BcryptPasswordHasher struct {
	cost int
}

// NewBcryptPasswordHasher creates a new bcrypt password hasher
func NewBcryptPasswordHasher() service.PasswordHasher {
	return &BcryptPasswordHasher{
		cost: bcrypt.DefaultCost,
	}
}

// Hash hashes a plain text password using bcrypt
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare compares a plain text password with a bcrypt hash
func (h *BcryptPasswordHasher) Compare(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
