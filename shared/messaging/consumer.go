package messaging

import (
	"encoding/json"
	"fmt"
	"log"
)

// EventConsumer consumes events from RabbitMQ
type EventConsumer struct {
	conn         *RabbitMQConnection
	exchangeName string
	queueName    string
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(conn *RabbitMQConnection, exchangeName, queueName, routingKey string) (*EventConsumer, error) {
	// Declare exchange
	if err := conn.DeclareExchange(exchangeName); err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue and bind to exchange
	if err := conn.DeclareQueue(queueName, exchangeName, routingKey); err != nil {
		return nil, fmt.Errorf("failed to declare/bind queue: %w", err)
	}

	return &EventConsumer{
		conn:         conn,
		exchangeName: exchangeName,
		queueName:    queueName,
	}, nil
}

// EventHandler is a function that handles consumed events
type EventHandler func(eventType string, body []byte) error

// Consume starts consuming messages from the queue
func (c *EventConsumer) Consume(handler EventHandler) error {
	msgs, err := c.conn.GetChannel().Consume(
		c.queueName, // queue
		"",          // consumer tag
		false,       // auto-ack (we'll manually ack)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started consuming from queue: %s", c.queueName)

	// Process messages
	go func() {
		for msg := range msgs {
			// Extract event type from message
			var eventTypeWrapper struct {
				EventType string `json:"event_type"`
			}

			if err := json.Unmarshal(msg.Body, &eventTypeWrapper); err != nil {
				log.Printf("Failed to unmarshal event type: %v", err)
				msg.Nack(false, false) // Reject message
				continue
			}

			// Handle event
			if err := handler(eventTypeWrapper.EventType, msg.Body); err != nil {
				log.Printf("Failed to handle event: %v", err)
				msg.Nack(false, true) // Requeue message
			} else {
				msg.Ack(false) // Acknowledge successful processing
			}
		}
	}()

	return nil
}

// Subscribe to a specific event type with a handler
func (c *EventConsumer) Subscribe(eventType string, handler func([]byte) error) error {
	// Wrap the handler to filter by event type
	eventHandler := func(receivedEventType string, body []byte) error {
		if receivedEventType == eventType {
			return handler(body)
		}
		return nil // Ignore events we're not interested in
	}

	return c.Consume(eventHandler)
}
