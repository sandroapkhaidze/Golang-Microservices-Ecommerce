package entity

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func (o *OrderItem) GetSubtotal() float64 {
	return float64(o.Quantity) * o.Price
}
