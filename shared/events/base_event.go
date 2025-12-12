package events

import (
    "time"
    "github.com/google/uuid"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
    EventID       string    `json:"event_id"`
    EventType     string    `json:"event_type"`
    AggregateID   string    `json:"aggregate_id"`
    OccurredAt    time.Time `json:"occurred_at"`
    CorrelationID string    `json:"correlation_id"` // For tracing across services
}

// NewBaseEvent creates a new base event with generated ID and timestamp
func NewBaseEvent(eventType, aggregateID, correlationID string) BaseEvent {
    return BaseEvent{
        EventID:       uuid.New().String(),
        EventType:     eventType,
        AggregateID:   aggregateID,
        OccurredAt:    time.Now().UTC(),
        CorrelationID: correlationID,
    }
}