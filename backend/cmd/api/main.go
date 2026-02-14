package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	admincmd "github.com/katerji/butchery-app/backend/internal/application/admin/commands"
	authcmd "github.com/katerji/butchery-app/backend/internal/application/auth/commands"
	custcmd "github.com/katerji/butchery-app/backend/internal/application/customer/commands"
	infraauth "github.com/katerji/butchery-app/backend/internal/infrastructure/auth"
	"github.com/katerji/butchery-app/backend/internal/infrastructure/persistence/postgres"
	apphttp "github.com/katerji/butchery-app/backend/internal/interface/http"
	"github.com/katerji/butchery-app/backend/internal/interface/http/handler"
	"github.com/katerji/butchery-app/backend/internal/interface/http/middleware"
	"github.com/katerji/butchery-app/backend/pkg/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Database
	ctx := context.Background()
	pool, err := postgres.NewConnectionPool(ctx, cfg.DB.DSN())
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	// Repositories
	adminRepo := postgres.NewAdminRepository(pool)
	customerRepo := postgres.NewCustomerRepository(pool)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(pool)

	// Infrastructure services
	passwordHasher := infraauth.NewBcryptHasher()
	tokenService := infraauth.NewTokenService(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)

	// Use case handlers
	adminLoginHandler := admincmd.NewAdminLoginHandler(adminRepo, passwordHasher, tokenService, refreshTokenRepo, cfg.JWT.AccessTokenTTL)
	registerCustomerHandler := custcmd.NewRegisterCustomerHandler(customerRepo, passwordHasher)
	customerLoginHandler := custcmd.NewCustomerLoginHandler(customerRepo, passwordHasher, tokenService, refreshTokenRepo, cfg.JWT.AccessTokenTTL)
	refreshTokenHandler := authcmd.NewRefreshTokenHandler(refreshTokenRepo, tokenService, cfg.JWT.AccessTokenTTL)
	logoutHandler := authcmd.NewLogoutHandler(refreshTokenRepo)

	// HTTP handlers
	adminAuthHandler := handler.NewAdminAuthHandler(adminLoginHandler)
	customerAuthHandler := handler.NewCustomerAuthHandler(registerCustomerHandler, customerLoginHandler)
	authHandler := handler.NewAuthHandler(refreshTokenHandler, logoutHandler)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// Router
	router := apphttp.NewRouter(apphttp.RouterDeps{
		Logger:              logger,
		AuthMiddleware:      authMiddleware,
		AdminAuthHandler:    adminAuthHandler,
		CustomerAuthHandler: customerAuthHandler,
		AuthHandler:         authHandler,
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("starting server", slog.String("addr", addr))

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
