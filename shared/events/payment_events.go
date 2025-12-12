package events

// Payment-related events

// PaymentProcessedEvent is published when payment is successful
type PaymentProcessedEvent struct {
    BaseEvent
    OrderID       string  `json:"order_id"`
    PaymentID     string  `json:"payment_id"`
    Amount        float64 `json:"amount"`
    PaymentMethod string  `json:"payment_method"`
}

// PaymentFailedEvent is published when payment fails
type PaymentFailedEvent struct {
    BaseEvent
    OrderID string `json:"order_id"`
    Reason  string `json:"reason"`
}

// Event type constants
const (
    PaymentProcessedEventType = "payment.processed"
    PaymentFailedEventType    = "payment.failed"
)