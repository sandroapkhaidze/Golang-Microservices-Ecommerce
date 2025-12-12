package events

// Inventory-related events

// InventoryReservedEvent is published when inventory is successfully reserved
type InventoryReservedEvent struct {
    BaseEvent
    OrderID     string                `json:"order_id"`
    Reservations []InventoryReservation `json:"reservations"`
}

// InventoryReservation represents a single item reservation
type InventoryReservation struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

// InventoryReservationFailedEvent is published when inventory reservation fails
type InventoryReservationFailedEvent struct {
    BaseEvent
    OrderID   string `json:"order_id"`
    ProductID string `json:"product_id"`
    Reason    string `json:"reason"`
}

// Event type constants
const (
    InventoryReservedEventType           = "inventory.reserved"
    InventoryReservationFailedEventType  = "inventory.reservation_failed"
)