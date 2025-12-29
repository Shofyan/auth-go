package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"auth-go/internal/application/dto"
	"auth-go/internal/application/usecase"
	"auth-go/internal/domain/repository"
	"auth-go/internal/interface/http/middleware"

	"github.com/google/uuid"
)

// WebHandler handles web UI requests
type WebHandler struct {
	templates           *template.Template
	logoutUseCase       *usecase.LogoutUseCase
	refreshTokenUseCase *usecase.RefreshTokenUseCase
	userRepo            repository.UserRepository
}

// NewWebHandler creates a new web handler
func NewWebHandler(
	logoutUseCase *usecase.LogoutUseCase,
	refreshTokenUseCase *usecase.RefreshTokenUseCase,
	userRepo repository.UserRepository,
) *WebHandler {
	// Parse all templates
	templates := template.Must(template.ParseGlob(filepath.Join("web", "templates", "*.html")))

	return &WebHandler{
		templates:           templates,
		logoutUseCase:       logoutUseCase,
		refreshTokenUseCase: refreshTokenUseCase,
		userRepo:            userRepo,
	}
}

// ServeLogin serves the login page
func (h *WebHandler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login",
	}
	// Parse login template with layout
	t := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "layout.html"),
		filepath.Join("web", "templates", "login.html"),
	))
	t.ExecuteTemplate(w, "layout.html", data)
}

// ServeRegister serves the register page
func (h *WebHandler) ServeRegister(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Register",
	}
	// Parse register template with layout
	t := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "layout.html"),
		filepath.Join("web", "templates", "register.html"),
	))
	t.ExecuteTemplate(w, "layout.html", data)
}

// ServeDashboard serves the dashboard page
func (h *WebHandler) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":     "Dashboard",
		"Dashboard": true,
	}
	// Parse dashboard template with layout
	t := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "layout.html"),
		filepath.Join("web", "templates", "dashboard.html"),
	))
	t.ExecuteTemplate(w, "layout.html", data)
}

// ServeProfile serves the profile page
func (h *WebHandler) ServeProfile(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":     "Profile",
		"Dashboard": true,
	}
	// Parse profile template with layout
	t := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "layout.html"),
		filepath.Join("web", "templates", "profile.html"),
	))
	t.ExecuteTemplate(w, "layout.html", data)
}

// ServeProfileData serves the profile data as HTML fragment
func (h *WebHandler) ServeProfileData(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<div class="error">Unauthorized</div>`))
		return
	}

	// Fetch full user data from repository
	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<div class="error">Failed to load user data</div>`))
		return
	}

	roleStrings := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleStrings[i] = role.String()
	}

	// Format dates
	createdAt := user.CreatedAt.Format("January 2, 2006 at 3:04 PM")
	updatedAt := user.UpdatedAt.Format("January 2, 2006 at 3:04 PM")
	lastLogin := "Never"
	if user.LastLoginAt != nil {
		lastLogin = user.LastLoginAt.Format("January 2, 2006 at 3:04 PM")
	}

	accountStatus := "Active"
	accountStatusColor := "#48bb78"
	if !user.IsActive {
		accountStatus = "Inactive"
		accountStatusColor = "#e53e3e"
	}

	// Return HTML fragment
	html := `<div class="profile-card">
		<div class="profile-field">
			<strong>User ID:</strong>
			<span style="font-family: monospace;">` + user.ID.String() + `</span>
		</div>
		<div class="profile-field">
			<strong>Email:</strong>
			<span>` + user.Email + `</span>
		</div>
		<div class="profile-field">
			<strong>Account Status:</strong>
			<span style="color: ` + accountStatusColor + `; font-weight: bold;">` + accountStatus + `</span>
		</div>
		<div class="profile-field">
			<strong>Roles:</strong>
			<div>`

	for _, role := range roleStrings {
		html += `<span class="badge">` + role + `</span>`
	}

	html += `</div>
		</div>
		<div class="profile-field">
			<strong>Account Created:</strong>
			<span>` + createdAt + `</span>
		</div>
		<div class="profile-field">
			<strong>Last Updated:</strong>
			<span>` + updatedAt + `</span>
		</div>
		<div class="profile-field">
			<strong>Last Login:</strong>
			<span>` + lastLogin + `</span>
		</div>
	</div>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// HandleLogout handles logout from web UI
func (h *WebHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := h.logoutUseCase.Execute(r.Context(), userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleRefreshToken handles token refresh from web UI
func (h *WebHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Try to get refresh token from request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`<div class="error">Invalid request</div>`))
		return
	}

	// Use the refresh token use case
	response, err := h.refreshTokenUseCase.Execute(r.Context(), dto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<div class="error">Token refresh failed</div>`))
		return
	}

	// Return success with new tokens as JSON for JavaScript to handle
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ServeHome redirects to login or dashboard
func (h *WebHandler) ServeHome(w http.ResponseWriter, r *http.Request) {
	// Serve the home/demo page
	http.ServeFile(w, r, filepath.Join("web", "templates", "home.html"))
}
