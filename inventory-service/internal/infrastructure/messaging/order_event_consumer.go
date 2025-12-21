package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/application/usecase"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/shared/events"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/shared/messaging"
)

type OrderEventConsumer struct {
	consumer            *messaging.EventConsumer
	reserveStockUseCase *usecase.ReserveStockUseCase
}

func NewOrderEventConsumer(
	consumer *messaging.EventConsumer,
	reserveStockUseCase *usecase.ReserveStockUseCase,
) *OrderEventConsumer {
	return &OrderEventConsumer{
		consumer:            consumer,
		reserveStockUseCase: reserveStockUseCase,
	}
}

// Start begins consuming order events
func (c *OrderEventConsumer) Start() error {
	log.Println("Starting Order Event Consumer...")

	// Subscribe to order.created events
	return c.consumer.Subscribe("order.created", c.handleOrderCreated)
}

// handleOrderCreated processes order.created events
func (c *OrderEventConsumer) handleOrderCreated(body []byte) error {
	log.Printf("Received order.created event: %s", string(body))

	// Unmarshal the event
	var event events.OrderCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("ERROR: Failed to unmarshal order.created event: %v", err)
		return err
	}

	log.Printf("Processing order %s (correlation_id: %s) with %d items",
		event.OrderID, event.CorrelationID, len(event.Items))

	// Execute the reserve stock use case
	ctx := context.Background()
	if err := c.reserveStockUseCase.Execute(ctx, event); err != nil {
		log.Printf("ERROR: Failed to reserve stock for order %s: %v", event.OrderID, err)
		// Error already handled in use case (failure event published)
		return err
	}

	log.Printf("Successfully reserved stock for order %s", event.OrderID)
	return nil
}
