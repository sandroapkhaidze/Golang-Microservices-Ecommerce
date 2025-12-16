package dto

import "time"

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	UserID string             `json:"user_id" binding:"required"`
	Items  []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

// OrderItemRequest represents a single item in the order request
type OrderItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,gt=0"`
}

// OrderResponse represents the order data returned to the client
type OrderResponse struct {
	ID            string              `json:"id"`
	UserID        string              `json:"user_id"`
	Status        string              `json:"status"`
	TotalAmount   float64             `json:"total_amount"`
	Items         []OrderItemResponse `json:"items"`
	CorrelationID string              `json:"correlation_id"`
	CreatedAt     time.Time           `json:"created_at"`
}

// OrderItemResponse represents a single item in the order response
type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
