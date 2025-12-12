package messaging

import (
    "fmt"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConnection manages the connection to RabbitMQ
type RabbitMQConnection struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    url     string
}

// NewRabbitMQConnection creates a new RabbitMQ connection
func NewRabbitMQConnection(host, port, user, password string) (*RabbitMQConnection, error) {
    url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
    
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    channel, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to open channel: %w", err)
    }

    log.Println("Successfully connected to RabbitMQ")

    return &RabbitMQConnection{
        conn:    conn,
        channel: channel,
        url:     url,
    }, nil
}

// Close closes the connection and channel
func (r *RabbitMQConnection) Close() error {
    if r.channel != nil {
        if err := r.channel.Close(); err != nil {
            return err
        }
    }
    if r.conn != nil {
        if err := r.conn.Close(); err != nil {
            return err
        }
    }
    return nil
}

// GetChannel returns the RabbitMQ channel
func (r *RabbitMQConnection) GetChannel() *amqp.Channel {
    return r.channel
}

// DeclareExchange declares a topic exchange
func (r *RabbitMQConnection) DeclareExchange(exchangeName string) error {
    return r.channel.ExchangeDeclare(
        exchangeName, // name
        "topic",      // type (topic allows routing by pattern)
        true,         // durable (survives server restart)
        false,        // auto-deleted
        false,        // internal
        false,        // no-wait
        nil,          // arguments
    )
}

// DeclareQueue declares a queue and binds it to an exchange
func (r *RabbitMQConnection) DeclareQueue(queueName, exchangeName, routingKey string) error {
    // Declare queue
    _, err := r.channel.QueueDeclare(
        queueName, // name
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare queue: %w", err)
    }

    // Bind queue to exchange
    err = r.channel.QueueBind(
        queueName,    // queue name
        routingKey,   // routing key
        exchangeName, // exchange
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to bind queue: %w", err)
    }

    return nil
}