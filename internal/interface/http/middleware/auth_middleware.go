package middleware

import (
	"context"
	"net/http"
	"strings"

	"auth-go/internal/domain/entity"
	"auth-go/internal/domain/service"
	apperrors "auth-go/pkg/errors"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UserRolesKey contextKey = "user_roles"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	tokenService service.TokenService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(tokenService service.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

// Authenticate validates JWT token and adds user context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.tokenService.ValidateAccessToken(token)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, apperrors.ErrInvalidToken.Error())
			return
		}

		// Add claims to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserRolesKey, claims.Roles)

		// Continue with authenticated request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole checks if user has required role (RBAC)
func (m *AuthMiddleware) RequireRole(role entity.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := r.Context().Value(UserRolesKey).([]entity.Role)
			if !ok {
				respondWithError(w, http.StatusForbidden, apperrors.ErrForbidden.Error())
				return
			}

			// Check if user has required role or higher
			hasPermission := false
			for _, userRole := range roles {
				if userRole >= role {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				respondWithError(w, http.StatusForbidden, apperrors.ErrForbidden.Error())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
