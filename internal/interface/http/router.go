package http

import (
	"net/http"

	"auth-go/internal/domain/entity"
	"auth-go/internal/interface/http/handler"
	"auth-go/internal/interface/http/middleware"
)

// Router sets up HTTP routes
type Router struct {
	authHandler    *handler.AuthHandler
	adminHandler   *handler.AdminHandler
	webHandler     *handler.WebHandler
	authMiddleware *middleware.AuthMiddleware
	logMiddleware  *middleware.LoggingMiddleware
	corsMiddleware *middleware.CORSMiddleware
}

// NewRouter creates a new router
func NewRouter(
	authHandler *handler.AuthHandler,
	adminHandler *handler.AdminHandler,
	webHandler *handler.WebHandler,
	authMiddleware *middleware.AuthMiddleware,
	logMiddleware *middleware.LoggingMiddleware,
	corsMiddleware *middleware.CORSMiddleware,
) *Router {
	return &Router{
		authHandler:    authHandler,
		adminHandler:   adminHandler,
		webHandler:     webHandler,
		authMiddleware: authMiddleware,
		logMiddleware:  logMiddleware,
		corsMiddleware: corsMiddleware,
	}
}

// Setup sets up all routes
func (rt *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	// API Routes
	// Public routes
	mux.HandleFunc("/api/v1/auth/register", rt.authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", rt.authHandler.Login)
	mux.HandleFunc("/api/v1/auth/refresh", rt.authHandler.RefreshToken)

	// Protected routes
	mux.Handle("/api/v1/auth/logout", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.authHandler.Logout)))
	mux.Handle("/api/v1/auth/profile", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.authHandler.GetProfile)))

	// Admin-only route example (RBAC)
	mux.Handle("/api/v1/admin/users",
		rt.authMiddleware.Authenticate(
			rt.authMiddleware.RequireRole(entity.RoleAdmin)(
				http.HandlerFunc(rt.adminHandler.ListUsers),
			),
		),
	)

	// Web UI Routes
	// Public web pages (HTML pages - authentication handled by JavaScript)
	mux.HandleFunc("/", rt.webHandler.ServeHome)
	mux.HandleFunc("/web/login", rt.webHandler.ServeLogin)
	mux.HandleFunc("/web/register", rt.webHandler.ServeRegister)
	mux.HandleFunc("/web/dashboard", rt.webHandler.ServeDashboard)
	mux.HandleFunc("/web/profile", rt.webHandler.ServeProfile)

	// Protected web data endpoints (API calls from JavaScript)
	mux.Handle("/web/profile-data", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.webHandler.ServeProfileData)))
	mux.Handle("/web/logout", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.webHandler.HandleLogout)))
	mux.Handle("/web/refresh-token", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.webHandler.HandleRefreshToken)))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Apply global middleware
	handler := rt.corsMiddleware.Handle(mux)
	handler = rt.logMiddleware.Log(handler)

	return handler
}
