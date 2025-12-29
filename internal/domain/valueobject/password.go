package valueobject

import (
	"errors"
	"unicode"
)

// Password represents a password value object
type Password struct {
	value string
}

// NewPassword creates a new password value object with validation
func NewPassword(password string) (*Password, error) {
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return nil, errors.New("password must not exceed 128 characters")
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return nil, errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return &Password{value: password}, nil
}

// Value returns the password value
func (p *Password) Value() string {
	return p.value
}
