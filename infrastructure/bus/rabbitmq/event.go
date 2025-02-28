package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"

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
	types    map[event.Type]reflect.Type
}

// RegisterEventType at startup
func (b *EventBus) RegisterEventType(eventType event.Type, eventStruct interface{}) {
	b.types[eventType] = reflect.TypeOf(eventStruct)
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
		log.Printf("Failed to declare exchange: %v", err)
		return nil, err
	}
	log.Printf("Exchange %s declared", exchange)

	return &EventBus{
		conn:     conn,
		channel:  channel,
		exchange: exchange,
		handlers: make(map[event.Type][]event.Handler),
		types:    make(map[event.Type]reflect.Type),
	}, nil
}

// Publish sends events to RabbitMQ.
func (b *EventBus) Publish(ctx context.Context, events []event.Event) error {
	log.Printf("Publishing %d events\n", len(events))
	if b.channel == nil {
		return errors.New("RabbitMQ channel is not open")
	}
	if b.conn == nil || b.conn.IsClosed() {
		return errors.New("RabbitMQ connection is not open")
	}
	for _, evt := range events {
		marshalledEvent, err := json.Marshal(evt)
		if err != nil {
			return err
		}
		envelope := event.EventEnvelope{
			EventType: evt.Type(),
			Data:      marshalledEvent,
		}
		data, err := json.Marshal(envelope)
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

// Consume listens for messages from the specified queue and dispatches them.
func (b *EventBus) Consume(queue string) error {
	log.Printf("Consuming from queue %s\n", queue)
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
		return err
	}
	log.Printf("Consumer successfully started on queue: %s", queue)
	for msg := range msgs { // Blocking loop, processes messages one by one
		log.Printf("Received message: %s", msg.Body)
		var envelope event.EventEnvelope
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			log.Printf("Failed to decode event envelope from queue %s: %v", queue, err)
			_ = msg.Nack(false, false) // Reject the message without requeueing
			continue
		}

		eventType, found := b.types[envelope.EventType]
		if !found {
			log.Printf("Unknown event type: %s", envelope.EventType)
			_ = msg.Nack(false, false)
			continue
		}

		evtValue := reflect.New(eventType).Interface()
		// Unmarshal into the correct event type
		if err := json.Unmarshal(envelope.Data, evtValue); err != nil {
			log.Printf("Failed to deserialize event: %v", err)
			_ = msg.Nack(false, false)
			continue
		}

		// Cast to event.Event interface
		evt, ok := evtValue.(event.Event)
		if !ok {
			log.Printf("Invalid event type: %T", evtValue)
			_ = msg.Nack(false, false)
			continue
		}

		handlers, ok := b.handlers[envelope.EventType]
		if !ok {
			log.Printf("No handlers for event type %s in queue %s", envelope.EventType, queue)
			_ = msg.Nack(false, false) // Reject the message without requeueing
			continue
		}

		// Process handlers synchronously (one at a time)
		for _, handler := range handlers {
			err := handler.Handle(context.Background(), evt)
			if err != nil {
				log.Printf("Error handling event %s from queue %s: %v", envelope.EventType, queue, err)
				_ = msg.Nack(false, true) // Requeue the message on failure
				continue
			}
		}
		log.Printf("Event %s processed", envelope.EventType)
		// Acknowledge the message after processing all handlers
		err = msg.Ack(false)
		if err != nil {
			log.Printf("Failed to acknowledge message from queue %s: %v", queue, err)
		} else {
			log.Printf("Message from queue %s acknowledged", queue)
		}
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
