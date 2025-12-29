package handler

import (
	"encoding/json"
	"net/http"

	"auth-go/internal/application/dto"
	"auth-go/internal/application/usecase"
	"auth-go/internal/interface/http/middleware"
	apperrors "auth-go/pkg/errors"

	"github.com/google/uuid"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	registerUseCase     *usecase.RegisterUseCase
	loginUseCase        *usecase.LoginUseCase
	refreshTokenUseCase *usecase.RefreshTokenUseCase
	logoutUseCase       *usecase.LogoutUseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	registerUseCase *usecase.RegisterUseCase,
	loginUseCase *usecase.LoginUseCase,
	refreshTokenUseCase *usecase.RefreshTokenUseCase,
	logoutUseCase *usecase.LogoutUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase:     registerUseCase,
		loginUseCase:        loginUseCase,
		refreshTokenUseCase: refreshTokenUseCase,
		logoutUseCase:       logoutUseCase,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.registerUseCase.Execute(r.Context(), req); err != nil {
		switch err {
		case apperrors.ErrInvalidEmail, apperrors.ErrInvalidPassword:
			respondWithError(w, http.StatusBadRequest, err.Error())
		case apperrors.ErrUserAlreadyExists:
			respondWithError(w, http.StatusConflict, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "user registered successfully"})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.loginUseCase.Execute(r.Context(), req)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidCredentials:
			respondWithError(w, http.StatusUnauthorized, err.Error())
		case apperrors.ErrUserInactive:
			respondWithError(w, http.StatusForbidden, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.refreshTokenUseCase.Execute(r.Context(), req)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidToken:
			respondWithError(w, http.StatusUnauthorized, err.Error())
		case apperrors.ErrExpiredToken:
			respondWithError(w, http.StatusUnauthorized, err.Error())
		case apperrors.ErrTokenReuse:
			respondWithError(w, http.StatusUnauthorized, "token reuse detected - all tokens revoked")
		case apperrors.ErrUserInactive:
			respondWithError(w, http.StatusForbidden, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.logoutUseCase.Execute(r.Context(), userID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// GetProfile returns the authenticated user's profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	email, _ := r.Context().Value(middleware.UserEmailKey).(string)
	roles, _ := r.Context().Value(middleware.UserRolesKey).([]interface{})

	roleStrings := make([]string, len(roles))
	for i, role := range roles {
		roleStrings[i] = role.(string)
	}

	response := dto.UserResponse{
		ID:    userID.String(),
		Email: email,
		Roles: roleStrings,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
