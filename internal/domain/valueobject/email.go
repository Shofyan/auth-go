package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Email represents an email value object
type Email struct {
	value string
}

// NewEmail creates a new email value object with validation
func NewEmail(email string) (*Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	return &Email{value: email}, nil
}

// Value returns the email value
func (e *Email) Value() string {
	return e.value
}

// String implements Stringer interface
func (e *Email) String() string {
	return e.value
}
