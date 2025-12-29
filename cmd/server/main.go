package main

import (
	"fmt"
	"log"
	"net/http"

	"auth-go/internal/application/usecase"
	"auth-go/internal/infrastructure/config"
	"auth-go/internal/infrastructure/persistence"
	"auth-go/internal/infrastructure/security"
	httpHandler "auth-go/internal/interface/http"
	"auth-go/internal/interface/http/handler"
	"auth-go/internal/interface/http/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := persistence.NewPostgresDB(persistence.DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Initialize repositories
	userRepo := persistence.NewPostgresUserRepository(db)
	refreshTokenRepo := persistence.NewPostgresRefreshTokenRepository(db)

	// Initialize services
	passwordHasher := security.NewBcryptPasswordHasher()
	tokenService := security.NewJWTTokenService(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
		cfg.JWT.Issuer,
	)

	// Initialize use cases
	registerUseCase := usecase.NewRegisterUseCase(userRepo, passwordHasher)
	loginUseCase := usecase.NewLoginUseCase(userRepo, refreshTokenRepo, passwordHasher, tokenService)
	refreshTokenUseCase := usecase.NewRefreshTokenUseCase(userRepo, refreshTokenRepo, tokenService)
	logoutUseCase := usecase.NewLogoutUseCase(refreshTokenRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(registerUseCase, loginUseCase, refreshTokenUseCase, logoutUseCase)
	adminHandler := handler.NewAdminHandler(userRepo)
	webHandler := handler.NewWebHandler(logoutUseCase, refreshTokenUseCase, userRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	logMiddleware := middleware.NewLoggingMiddleware()
	corsMiddleware := middleware.NewCORSMiddleware()

	// Setup router
	router := httpHandler.NewRouter(authHandler, adminHandler, webHandler, authMiddleware, logMiddleware, corsMiddleware)
	httpHandler := router.Setup()

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)

	if err := http.ListenAndServe(addr, httpHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
