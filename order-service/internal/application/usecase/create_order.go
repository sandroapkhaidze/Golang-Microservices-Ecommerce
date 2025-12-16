package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/application/dto"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/domain/entity"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/domain/repository"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/shared/events"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/shared/messaging"
)

type CreateOrderUseCase struct {
	orderRepo      repository.OrderRepository
	eventPublisher *messaging.EventPublisher
}

// NewCreateOrderUseCase creates a new CreateOrderUseCase
func NewCreateOrderUseCase(
	orderRepo repository.OrderRepository,
	eventPublisher *messaging.EventPublisher,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// 1. Generate IDs
	orderID := uuid.New().String()
	correlationID := uuid.New().String()

	// 2. Convert request items to entity items
	items := make([]entity.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = entity.OrderItem{
			ID:        uuid.New().String(),
			OrderID:   orderID,        // ← Fixed
			ProductID: item.ProductID, // ← Fixed
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// 3. Create order entity
	order := &entity.Order{
		ID:            orderID,
		UserID:        req.UserID,
		Status:        entity.OrderStatusPending,
		Items:         items,
		CorrelationID: correlationID,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	// 4. Calculate total
	order.CalculateTotal()

	// 5. Save to database
	err := uc.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 6. Publish event
	event := events.OrderCreatedEvent{
		BaseEvent: events.NewBaseEvent(
			events.OrderCreatedEventType,
			orderID,
			correlationID,
		),
		OrderID:     orderID,
		UserID:      req.UserID,
		TotalAmount: order.TotalAmount,
		Items:       convertToEventItems(req.Items),
	}

	err = uc.eventPublisher.Publish("order.created", event)
	if err != nil {
		// Log but don't fail the order
		log.Printf("Warning: failed to publish OrderCreatedEvent: %v", err) // ← Fixed
	}

	// 7. Convert to response DTO
	return &dto.OrderResponse{
		ID:            order.ID,
		UserID:        order.UserID,
		Status:        string(order.Status),
		TotalAmount:   order.TotalAmount,
		Items:         convertToResponseItems(order.Items),
		CorrelationID: order.CorrelationID,
		CreatedAt:     order.CreatedAt,
	}, nil
}

func convertToResponseItems(items []entity.OrderItem) []dto.OrderItemResponse {
	orderItems := make([]dto.OrderItemResponse, len(items))
	for i, item := range items {
		orderItems[i] = dto.OrderItemResponse{
			ID:        item.ID, // ← Fixed
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}
	return orderItems
}

// Helper: Convert request items to event items
func convertToEventItems(items []dto.OrderItemRequest) []events.OrderItem {
	orderItems := make([]events.OrderItem, len(items))
	for i, item := range items {
		orderItems[i] = events.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}
	return orderItems
}
