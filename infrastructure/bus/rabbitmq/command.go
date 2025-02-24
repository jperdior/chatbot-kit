package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/jperdior/chatbot-kit/application/command"
	"github.com/streadway/amqp"
	"log"
)

type CommandBus struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queue    string
	handlers map[command.Type][]command.Handler
}

// NewCommandBus initializes a new RabbitMQ-based CommandBus.
func NewCommandBus(amqpURL, exchange, queue string) (*CommandBus, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	// Declare exchange (only relevant for publishing)
	err = ch.ExchangeDeclare(
		exchange,
		"direct",
		true,  // Durable
		false, // Auto-delete
		false, // Internal
		false, // No-wait
		nil,
	)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &CommandBus{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
		queue:    queue,
		handlers: make(map[command.Type][]command.Handler),
	}, nil
}

// Publish sends a command to the bus.
func (b *CommandBus) Publish(ctx context.Context, cmd command.Command) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}

	return b.channel.Publish(b.exchange, b.queue, false, false, msg)
}

// Subscribe registers a command handler.
func (b *CommandBus) Subscribe(cmdType command.Type, handler command.Handler) {
	b.handlers[cmdType] = append(b.handlers[cmdType], handler)
}

// Consume listens for messages from the queue and dispatches them.
func (b *CommandBus) Consume() error {
	_, err := b.channel.QueueDeclare(
		b.queue,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,
	)
	if err != nil {
		return err
	}

	if err := b.channel.QueueBind(b.queue, b.queue, b.exchange, false, nil); err != nil {
		return err
	}

	msgs, err := b.channel.Consume(
		b.queue,
		"",    // Consumer name
		false, // Auto-ack (false for manual ack)
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		var cmd command.Command
		if err := json.Unmarshal(msg.Body, &cmd); err != nil {
			log.Printf("Failed to decode command: %v", err)
			_ = msg.Nack(false, false) // Reject without requeueing
			continue
		}

		handlers, ok := b.handlers[cmd.Type()]
		if !ok {
			log.Printf("No handlers for command type: %s", cmd.Type())
			_ = msg.Nack(false, false)
			continue
		}

		for _, handler := range handlers {
			if err := handler.Handle(context.Background(), cmd); err != nil {
				log.Printf("Error handling command %s: %v", cmd.Type(), err)
				_ = msg.Nack(false, true) // Requeue on failure
				continue
			}
		}

		_ = msg.Ack(false) // Acknowledge message
	}

	return nil
}

// Close cleans up connections and channels.
func (b *CommandBus) Close() {
	if b.channel != nil {
		_ = b.channel.Close()
	}
	if b.conn != nil {
		_ = b.conn.Close()
	}
}
