package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user aggregate root in DDD
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Roles        []Role
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time
}

// NewUser creates a new user entity
func NewUser(email, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		Roles:        []Role{RoleUser}, // Default role
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// HasRole checks if user has a specific role (RBAC)
func (u *User) HasRole(role Role) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasPermission checks if user has permission based on role hierarchy
func (u *User) HasPermission(requiredRole Role) bool {
	if !u.IsActive {
		return false
	}

	for _, userRole := range u.Roles {
		if userRole >= requiredRole {
			return true
		}
	}
	return false
}

// AddRole adds a role to the user (RBAC)
func (u *User) AddRole(role Role) {
	if !u.HasRole(role) {
		u.Roles = append(u.Roles, role)
		u.UpdatedAt = time.Now()
	}
}

// RemoveRole removes a role from the user (RBAC)
func (u *User) RemoveRole(role Role) {
	for i, r := range u.Roles {
		if r == role {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			u.UpdatedAt = time.Now()
			break
		}
	}
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}
