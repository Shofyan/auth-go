package handler

import (
	"encoding/json"
	"html/template"
	"log"
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

	// Fetch full user data from repository (PostgreSQL database)
	log.Printf("Fetching user data from PostgreSQL database for user ID: %s", userID.String())
	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil {
		log.Printf("Error fetching user from database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<div class="error">Failed to load user data</div>`))
		return
	}
	log.Printf("Successfully retrieved user data from database: email=%s, roles=%v", user.Email, user.Roles)

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

	// Return HTML fragment with table format
	html := `
	<div class="profile-card">
		<h3 style="margin-top: 0; margin-bottom: 10px; color: #4a5568;">Account Information</h3>
		<p style="font-size: 12px; color: #718096; margin-bottom: 15px;">📊 Data retrieved from PostgreSQL database for User ID: ` + user.ID.String() + `</p>
		<table style="width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden;">
			<thead>
				<tr style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white;">
					<th style="padding: 12px; text-align: left; font-weight: 600;">Field</th>
					<th style="padding: 12px; text-align: left; font-weight: 600;">Value</th>
				</tr>
			</thead>
			<tbody>
				<tr style="border-bottom: 1px solid #e2e8f0;">
					<td style="padding: 12px; font-weight: 600; color: #4a5568;">User ID</td>
					<td style="padding: 12px; font-family: monospace; color: #2d3748;">` + user.ID.String() + `</td>
				</tr>
				<tr style="border-bottom: 1px solid #e2e8f0; background: #f7fafc;">
					<td style="padding: 12px; font-weight: 600; color: #4a5568;">Email</td>
					<td style="padding: 12px; color: #2d3748;">` + user.Email + `</td>
				</tr>
				<tr style="border-bottom: 1px solid #e2e8f0;">
					<td style="padding: 12px; font-weight: 600; color: #4a5568;">Account Status</td>
					<td style="padding: 12px;">
						<span style="color: ` + accountStatusColor + `; font-weight: bold; padding: 4px 12px; background: ` + accountStatusColor + `20; border-radius: 12px;">
							` + accountStatus + `
						</span>
					</td>
				</tr>
				<tr style="border-bottom: 1px solid #e2e8f0; background: #f7fafc;">
					<td style="padding: 12px; font-weight: 600; color: #4a5568;">Roles</td>
					<td style="padding: 12px;">`

	for _, role := range roleStrings {
		html += `<span class="badge" style="margin-right: 6px;">` + role + `</span>`
	}

	html += `</td>
				</tr>
				<tr style="border-bottom: 1px solid #e2e8f0;">
					<td style="padding: 12px; font-weight: 600; color: #4a5568;">Account Created</td>
					<td style="padding: 12px; color: #2d3748;">` + createdAt + `</td>
				</tr>
				<tr style="border-bottom: 1px solid #e2e8f0; background: #f7fafc;">
					<td style="padding: 12px; font-weight: 600; color: #4a5568;">Last Updated</td>
					<td style="padding: 12px; color: #2d3748;">` + updatedAt + `</td>
				</tr>
				<tr>
					<td style="padding: 12px; font-weight: 600; color: #4a5568;">Last Login</td>
					<td style="padding: 12px; color: #2d3748;">` + lastLogin + `</td>
				</tr>
			</tbody>
		</table>
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
