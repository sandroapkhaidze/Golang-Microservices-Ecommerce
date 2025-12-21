package usecase

import (
	"context"
	"fmt"

	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/domain/entity"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/domain/repository"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/shared/events"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/shared/messaging"
)

type ReserveStockUseCase struct {
	inventoryRepo  repository.InventoryRepository
	eventPublisher *messaging.EventPublisher
}

func NewReserveStockUseCase(
	inventoryRepo repository.InventoryRepository,
	eventPublisher *messaging.EventPublisher) *ReserveStockUseCase {
	return &ReserveStockUseCase{
		inventoryRepo:  inventoryRepo,
		eventPublisher: eventPublisher,
	}
}

func (uc *ReserveStockUseCase) Execute(ctx context.Context, event events.OrderCreatedEvent) error {
	productIDs := make([]string, len(event.Items))
	for _, item := range event.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := uc.inventoryRepo.GetByIDs(ctx, productIDs)
	if err != nil {
		uc.publishFailureEvent(event.CorrelationID, event.OrderID, "failed to fetch products")
		return fmt.Errorf("failed to get products: %w", err)
	}

	productMap := make(map[string]*entity.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	for _, item := range event.Items {
		product, exists := productMap[item.ProductID]

		// Check if product exists
		if !exists {
			msg := fmt.Sprintf("product %s not found", item.ProductID)
			uc.publishFailureEvent(event.CorrelationID, event.OrderID, msg)
			return fmt.Errorf(msg)
		}

		// Check if we can reserve (product active, enough stock, etc.)
		if err := product.CanReserve(int32(item.Quantity)); err != nil {
			msg := fmt.Sprintf("cannot reserve product %s: %v", item.ProductID, err)
			uc.publishFailureEvent(event.CorrelationID, event.OrderID, msg)
			return fmt.Errorf(msg)
		}
	}

	for _, item := range event.Items {
		product := productMap[item.ProductID]
		if err := product.ReserveStock(int32(item.Quantity)); err != nil {
			uc.publishFailureEvent(event.CorrelationID, event.OrderID, err.Error())
			return err
		}
	}

	if err := uc.inventoryRepo.UpdateMultiple(ctx, products); err != nil {
		msg := "failed to save reserved stock to database"
		uc.publishFailureEvent(event.CorrelationID, event.OrderID, msg)
		return fmt.Errorf("%s: %w", msg, err)
	}

	uc.publishSuccessEvent(event.CorrelationID, event.OrderID, event.Items)

	return nil

}

func (uc *ReserveStockUseCase) publishSuccessEvent(correlationID, orderID string, items []events.OrderItem) {
	reservations := make([]events.InventoryReservation, len(items))
	for i, item := range items {
		reservations[i] = events.InventoryReservation{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	reservedEvent := events.InventoryReservedEvent{
		BaseEvent: events.BaseEvent{
			CorrelationID: correlationID,
		},
		OrderID:      orderID,
		Reservations: reservations,
	}

	if err := uc.eventPublisher.Publish("inventory.reserved", reservedEvent); err != nil {
		// Log error but don't fail the operation (DB already saved)
		fmt.Printf("ERROR: Failed to publish inventory.reserved event: %v\n", err)
	}
}

func (uc *ReserveStockUseCase) publishFailureEvent(correlationID, orderID, reason string) {
	failedEvent := events.InventoryReservationFailedEvent{
		BaseEvent: events.BaseEvent{
			CorrelationID: correlationID,
		},
		OrderID: orderID,
		Reason:  reason,
	}

	if err := uc.eventPublisher.Publish("inventory.reservation_failed", failedEvent); err != nil {
		fmt.Printf("ERROR: Failed to publish inventory.reservation_failed event: %v\n", err)
	}
}
