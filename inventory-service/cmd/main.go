package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/application/usecase"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/infrastructure/config"
	infraMessaging "github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/infrastructure/messaging"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/infrastructure/persistence"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/shared/messaging"
)

func main() {
	log.Println("Starting Inventory Service...")

	// Initialize database
	db := config.NewDatabase()
	defer db.Close()

	// Initialize repository
	inventoryRepo := persistence.NewPostgresInventoryRepository(db)

	// Initialize RabbitMQ connection
	rabbitConn, err := messaging.NewRabbitMQConnection(
		getEnv("RABBITMQ_HOST", "localhost"),
		getEnv("RABBITMQ_PORT", "5672"),
		getEnv("RABBITMQ_USER", "admin"),
		getEnv("RABBITMQ_PASSWORD", "admin"),
	)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitConn.Close()

	// Initialize publisher
	publisher, err := messaging.NewEventPublisher(rabbitConn, "ecommerce-events")
	if err != nil {
		log.Fatal("Failed to create publisher:", err)
	}

	// Initialize consumer
	consumer, err := messaging.NewEventConsumer(
		rabbitConn,
		"ecommerce-events",  // exchange name
		"inventory-service", // queue name
		"order.created",     // routing key
	)
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}

	// Initialize use cases
	reserveStockUseCase := usecase.NewReserveStockUseCase(inventoryRepo, publisher)

	// Initialize event consumer
	orderEventConsumer := infraMessaging.NewOrderEventConsumer(consumer, reserveStockUseCase)

	// Start consuming events
	if err := orderEventConsumer.Start(); err != nil {
		log.Fatal("Failed to start event consumer:", err)
	}

	log.Println("✅ Inventory Service is running and listening for events...")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Inventory Service...")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
