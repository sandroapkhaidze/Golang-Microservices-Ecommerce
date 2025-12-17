package main

import (
	"log"
	"os"

	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/application/usecase"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/infrastructure/config"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/infrastructure/persistence"
	httpHandler "github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/presentation/http"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/shared/messaging"
)

func main() {
	// Load configuration from environment variables
	dbConfig := config.DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5433"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "orderdb"),
	}

	// Connect to database
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to RabbitMQ
	rabbitConn, err := messaging.NewRabbitMQConnection(
		getEnv("RABBITMQ_HOST", "localhost"),
		getEnv("RABBITMQ_PORT", "5672"),
		getEnv("RABBITMQ_USER", "admin"),
		getEnv("RABBITMQ_PASSWORD", "admin"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Create event publisher
	eventPublisher, err := messaging.NewEventPublisher(
		rabbitConn,
		"ecommerce-events",
	)
	if err != nil {
		log.Fatalf("Failed to create event publisher: %v", err)
	}

	// Initialize repository
	orderRepo := persistence.NewPostgresOrderRepository(db)

	// Initialize use case
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, eventPublisher)

	// Initialize HTTP handler
	orderHandler := httpHandler.NewOrderHandler(createOrderUseCase)

	// Setup router
	router := httpHandler.SetupRouter(orderHandler)

	// Start server
	port := getEnv("PORT", "8082")
	log.Printf("Order Service starting on port %s", port)
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
