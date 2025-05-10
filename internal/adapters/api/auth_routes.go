package api

import (
	"auth-service/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// SetupAuthRoutes sets up all authentication-related routes
func SetupAuthRoutes(r chi.Router, authHandler *AuthHandler) {
	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		// Add your protected routes here
		// Example:
		// r.Get("/api/profile", authHandler.GetProfile)
	})
}
