package handler

import (
	"net/http"

	"auth-go/internal/domain/repository"
)

// AdminHandler handles admin HTTP requests
type AdminHandler struct {
	userRepo repository.UserRepository
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(userRepo repository.UserRepository) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
	}
}

// ListUsers lists all users (admin only)
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.FindAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	// Convert to response format
	type UserResponse struct {
		ID        string   `json:"id"`
		Email     string   `json:"email"`
		Roles     []string `json:"roles"`
		IsActive  bool     `json:"is_active"`
		CreatedAt string   `json:"created_at"`
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		roles := make([]string, len(user.Roles))
		for j, role := range user.Roles {
			roles[j] = role.String()
		}

		response[i] = UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Roles:     roles,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}

// // Helper functions
// func respondWithError(w http.ResponseWriter, code int, message string) {
// 	respondWithJSON(w, code, map[string]string{"error": message})
// }

// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
// 	response, _ := json.Marshal(payload)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// }
