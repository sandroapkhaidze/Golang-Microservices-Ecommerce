package entity

import "time"

type Order struct {
	ID            string      `json:"id"`
	UserID        string      `json:"user_id"`
	Status        OrderStatus `json:"status"`
	TotalAmount   float64     `json:"total_amount"`
	Items         []OrderItem `json:"items"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	CorrelationID string      `json:"correlation_id"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusFailed     OrderStatus = "failed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

func (o *Order) CalculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	o.TotalAmount = total
}

func (o *Order) MarkAsProcessing() {
	o.Status = OrderStatusProcessing
	o.UpdatedAt = time.Now().UTC()
}

func (o *Order) MarkAsCompleted() {
	o.Status = OrderStatusCompleted
	o.UpdatedAt = time.Now().UTC()
}

func (o *Order) MarkAsFailed() {
	o.Status = OrderStatusFailed
	o.UpdatedAt = time.Now().UTC()
}

func (o *Order) MarkAsCancelled() {
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now().UTC()
}

func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusProcessing
}
