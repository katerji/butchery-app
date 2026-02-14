package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/katerji/butchery-app/backend/internal/interface/http/handler"
	"github.com/katerji/butchery-app/backend/internal/interface/http/middleware"
)

// RouterDeps holds all handler and middleware dependencies for the router.
type RouterDeps struct {
	Logger              *slog.Logger
	AuthMiddleware      *middleware.AuthMiddleware
	AdminAuthHandler    *handler.AdminAuthHandler
	CustomerAuthHandler *handler.CustomerAuthHandler
	AuthHandler         *handler.AuthHandler
}

// NewRouter creates a new chi router with all routes and middleware.
func NewRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(requestLogger(deps.Logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		// Public customer auth routes
		r.Post("/auth/register", deps.CustomerAuthHandler.Register)
		r.Post("/auth/login", deps.CustomerAuthHandler.Login)

		// Public token refresh
		r.Post("/auth/refresh", deps.AuthHandler.Refresh)

		// Public admin auth routes
		r.Post("/admin/auth/login", deps.AdminAuthHandler.Login)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(deps.AuthMiddleware.RequireAuth)
			r.Post("/auth/logout", deps.AuthHandler.Logout)
		})
	})

	return r
}

func requestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Duration("duration", time.Since(start)),
				slog.String("request_id", chimw.GetReqID(r.Context())),
			)
		})
	}
}
