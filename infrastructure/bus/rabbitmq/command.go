package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/jperdior/chatbot-kit/application/command"
	"github.com/streadway/amqp"
	"log"
)

type CommandBus struct {
	publishConn    *amqp.Connection
	consumeConn    *amqp.Connection
	pubChannel     *amqp.Channel
	consumeChannel *amqp.Channel
	exchange       string
	queue          string
	handlers       map[command.Type][]command.Handler
}

// NewCommandBus initializes a new RabbitMQ-based EventBus.
func NewCommandBus(amqpURL, exchange, queue string) (*CommandBus, error) {
	publishConn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	consumeConn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	pubCh, err := publishConn.Channel()
	if err != nil {
		return nil, err
	}

	consumeCh, err := consumeConn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare exchange
	err = pubCh.ExchangeDeclare(
		exchange,
		"fanout",
		true,  // Durable
		false, // Auto-delete
		false, // Internal
		false, // No-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Declare queue
	_, err = consumeCh.QueueDeclare(
		queue,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Bind queue to exchange
	if err = consumeCh.QueueBind(queue, "", exchange, false, nil); err != nil {
		return nil, err
	}

	return &CommandBus{
		publishConn:    publishConn,
		consumeConn:    consumeConn,
		pubChannel:     pubCh,
		consumeChannel: consumeCh,
		exchange:       exchange,
		queue:          queue,
		handlers:       make(map[command.Type][]command.Handler),
	}, nil
}

// Publish sends command to the bus.
func (b *CommandBus) Publish(ctx context.Context, cmd command.Command) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}

	if err := b.pubChannel.Publish(b.exchange, "", false, false, msg); err != nil {
		return err
	}
	return nil
}

// Subscribe registers an event handler.
func (b *CommandBus) Subscribe(cmdType command.Type, handler command.Handler) {
	b.handlers[cmdType] = append(b.handlers[cmdType], handler)
}

// Consume listens for messages from the queue and dispatches them.
func (b *CommandBus) Consume() error {
	msgs, err := b.consumeChannel.Consume(
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
			log.Printf("Failed to decode event: %v", err)
			_ = msg.Nack(false, false) // Reject the message without requeueing
			continue
		}

		handlers, ok := b.handlers[cmd.Type()]
		if !ok {
			log.Printf("No handlers for event type: %s", cmd.Type())
			_ = msg.Nack(false, false) // Reject the message without requeueing
			continue
		}

		// Process handlers synchronously (one at a time)
		for _, handler := range handlers {
			err := handler.Handle(context.Background(), cmd)
			if err != nil {
				log.Printf("Error handling event %s: %v", cmd.Type(), err)
				_ = msg.Nack(false, true) // Requeue the message on failure
				continue
			}
		}

		// Acknowledge the message after processing all handlers
		err = msg.Ack(false)
		if err != nil {
			log.Printf("Failed to acknowledge message: %v", err)
			return err
		}
	}

	return nil
}

// Close cleans up connections and channels.
func (b *CommandBus) Close() {
	if b.pubChannel != nil {
		_ = b.pubChannel.Close()
	}
	if b.consumeChannel != nil {
		_ = b.consumeChannel.Close()
	}
	if b.publishConn != nil {
		_ = b.publishConn.Close()
	}
	if b.consumeConn != nil {
		_ = b.consumeConn.Close()
	}
}
