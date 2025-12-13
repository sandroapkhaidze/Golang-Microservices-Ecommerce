package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EventPublisher publishes events to RabbitMQ
type EventPublisher struct {
	conn         *RabbitMQConnection
	exchangeName string
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(conn *RabbitMQConnection, exchangeName string) (*EventPublisher, error) {
	// Declare exchange
	if err := conn.DeclareExchange(exchangeName); err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &EventPublisher{
		conn:         conn,
		exchangeName: exchangeName,
	}, nil
}

// Publish publishes an event to RabbitMQ
func (p *EventPublisher) Publish(routingKey string, event interface{}) error {
	// Serialize event to JSON
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publish message
	err = p.conn.GetChannel().PublishWithContext(
		ctx,
		p.exchangeName, // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}
