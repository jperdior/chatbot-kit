package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jperdior/chatbot-kit/application/event"
	"github.com/streadway/amqp"
)

// EventBus is a RabbitMQ implementation of the event.Bus.
type EventBus struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queues   []string
	handlers map[event.Type][]event.Handler
}

// NewEventBus initializes a new RabbitMQ-based EventBus.
func NewEventBus(amqpURL, exchange string) (*EventBus, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare exchange
	err = channel.ExchangeDeclare(
		exchange,
		"topic", // Change to "topic" to support routing keys
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &EventBus{
		conn:     conn,
		channel:  channel,
		exchange: exchange,
		handlers: make(map[event.Type][]event.Handler),
	}, nil
}

// Publish sends events to RabbitMQ.
func (b *EventBus) Publish(ctx context.Context, events []event.Event) error {
	for _, evt := range events {
		data, err := json.Marshal(evt)
		if err != nil {
			return err
		}

		msg := amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		}

		routingKey := string(evt.Type())

		if err := b.channel.Publish(b.exchange, routingKey, false, false, msg); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe registers an event handler.
func (b *EventBus) Subscribe(evtType event.Type, handler event.Handler) {
	b.handlers[evtType] = append(b.handlers[evtType], handler)
}

// BindQueue binds a queue to the exchange with a routing key.
func (b *EventBus) BindQueue(queue, routingKey string) error {
	_, err := b.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = b.channel.QueueBind(queue, routingKey, b.exchange, false, nil)
	if err != nil {
		return err
	}
	b.queues = append(b.queues, queue)
	return nil
}

// Consume listens for messages from the queue and dispatches them.
// Consume listens for messages from the queue and dispatches them.
func (b *EventBus) Consume() error {
	for _, queue := range b.queues { // Iterate over all bound queues
		go func(queue string) { // Start a goroutine for each queue
			msgs, err := b.channel.Consume(
				queue,
				"",    // Consumer name
				false, // Auto-ack (false for manual ack)
				false, // Exclusive
				false, // No-local
				false, // No-wait
				nil,
			)
			if err != nil {
				log.Printf("Failed to start consuming from queue %s: %v", queue, err)
				return
			}

			for msg := range msgs { // Blocking loop, processes messages one by one
				var evt event.Event
				if err := json.Unmarshal(msg.Body, &evt); err != nil {
					log.Printf("Failed to decode event: %v", err)
					_ = msg.Nack(false, false) // Reject the message without requeueing
					continue
				}

				handlers, ok := b.handlers[evt.Type()]
				if !ok {
					log.Printf("No handlers for event type: %s", evt.Type())
					_ = msg.Nack(false, false) // Reject the message without requeueing
					continue
				}

				// Process handlers synchronously (one at a time)
				for _, handler := range handlers {
					err := handler.Handle(context.Background(), evt)
					if err != nil {
						log.Printf("Error handling event %s: %v", evt.Type(), err)
						_ = msg.Nack(false, true) // Requeue the message on failure
						continue
					}
				}

				// Acknowledge the message after processing all handlers
				err = msg.Ack(false)
				if err != nil {
					log.Printf("Failed to acknowledge message: %v", err)
				}
			}
		}(queue)
	}

	return nil
}

// Close cleans up connections and channels.
func (b *EventBus) Close() {
	if b.channel != nil {
		_ = b.channel.Close()
	}
	if b.conn != nil {
		_ = b.conn.Close()
	}
}
