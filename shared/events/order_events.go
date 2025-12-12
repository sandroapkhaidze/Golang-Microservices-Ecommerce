package events

// Order lifecycle events

// OrderCreatedEvent is published when a new order is created
type OrderCreatedEvent struct {
    BaseEvent
    OrderID      string  `json:"order_id"`
    UserID       string  `json:"user_id"`
    TotalAmount  float64 `json:"total_amount"`
    Items        []OrderItem `json:"items"`
}

// OrderItem represents an item in the order
type OrderItem struct {
    ProductID string  `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

// OrderCompletedEvent is published when order processing is successful
type OrderCompletedEvent struct {
    BaseEvent
    OrderID string `json:"order_id"`
    UserID  string `json:"user_id"`
}

// OrderFailedEvent is published when order processing fails
type OrderFailedEvent struct {
    BaseEvent
    OrderID string `json:"order_id"`
    UserID  string `json:"user_id"`
    Reason  string `json:"reason"`
}

// Event type constants
const (
    OrderCreatedEventType   = "order.created"
    OrderCompletedEventType = "order.completed"
    OrderFailedEventType    = "order.failed"
)
