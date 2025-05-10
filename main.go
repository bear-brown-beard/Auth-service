package main

import (
	"fmt"
	"log"
	"net/http"

	"auth-service/internal/adapters/api"
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/repositories"
	"auth-service/internal/services"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode,
	)

	err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Initialize services
	authRepo := repositories.NewAuthRepository(database.DB)
	authService := services.NewAuthService(authRepo, cfg.DBPassword)
	authHandler := api.NewAuthHandler(authService)

	// Create router
	r := chi.NewRouter()

	// Setup routes
	api.SetupAuthRoutes(r, authHandler)

	// Start server
	port := fmt.Sprintf(":%s", cfg.DBPort)
	log.Printf("Auth service is running on %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
