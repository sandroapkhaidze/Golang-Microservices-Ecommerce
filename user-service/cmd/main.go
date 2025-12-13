package main

import (
	"log"
	"os"

	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/application/usecase"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/infrastructure/config"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/infrastructure/persistence"
	httpHandler "github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/presentation/http"
)

func main() {
	// Load configuration from environment variables
	dbConfig := config.DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "userdb"),
	}

	// Connect to database
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	userRepo := persistence.NewPostgresUserRepository(db)

	// Initialize use cases
	registerUseCase := usecase.NewRegisterUserUseCase(userRepo)
	loginUseCase := usecase.NewLoginUserUseCase(userRepo)

	// Initialize HTTP handler
	userHandler := httpHandler.NewUserHandler(registerUseCase, loginUseCase)

	// Setup router
	router := httpHandler.SetupRouter(userHandler)

	// Start server
	port := getEnv("PORT", "8081")
	log.Printf("User Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
