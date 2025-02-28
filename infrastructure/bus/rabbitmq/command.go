package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/jperdior/chatbot-kit/application/command"
	"github.com/streadway/amqp"
	"log"
	"reflect"
)

type CommandBus struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queue    string
	handlers map[command.Type][]command.Handler
	types    map[command.Type]reflect.Type
}

// RegisterCommandType at startup
func (b *CommandBus) RegisterCommandType(commandType command.Type, commandStruct interface{}) {
	b.types[commandType] = reflect.TypeOf(commandStruct)
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
		types:    make(map[command.Type]reflect.Type),
	}, nil
}

// Dispatch sends a command to the bus.
func (b *CommandBus) Dispatch(ctx context.Context, cmd command.Command) error {
	log.Printf("Dispatching command: %s", cmd.Type())
	marshalledCommand, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Failed to marshal command: %v", err)
		return err
	}

	envelope := command.CommandEnvelope{
		CommandType: cmd.Type(),
		Data:        marshalledCommand,
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}

	err = b.channel.Publish(b.exchange, b.queue, false, false, msg)
	if err != nil {
		log.Printf("Failed to publish command: %v", err)
		return err
	}

	log.Printf("Command dispatched: %s", cmd.Type())
	return nil
}

// Register registers a command handler.
func (b *CommandBus) Register(cmdType command.Type, handler command.Handler) {
	b.handlers[cmdType] = append(b.handlers[cmdType], handler)
}

// Consume listens for messages from the queue and dispatches them.
func (b *CommandBus) Consume() error {
	log.Printf("Starting to consume from queue: %s", b.queue)
	_, err := b.channel.QueueDeclare(
		b.queue,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return err
	}

	if err := b.channel.QueueBind(b.queue, b.queue, b.exchange, false, nil); err != nil {
		log.Printf("Failed to bind queue: %v", err)
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
		log.Printf("Failed to start consuming: %v", err)
		return err
	}

	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
		var envelope command.CommandEnvelope
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			log.Printf("Failed to decode event envelope from queue %s: %v", b.queue, err)
			_ = msg.Nack(false, false) // Reject the message without requeueing
			continue
		}
		commandType, found := b.types[envelope.CommandType]
		if !found {
			log.Printf("Unknown event type: %s", envelope.CommandType)
			_ = msg.Nack(false, false)
			continue
		}

		commandValue := reflect.New(commandType).Interface()
		if err := json.Unmarshal(envelope.Data, commandValue); err != nil {
			log.Printf("Failed to deserialize event: %v", err)
			_ = msg.Nack(false, false)
			continue
		}

		cmd, ok := commandValue.(command.Command)
		if !ok {
			log.Printf("Invalid event type: %T", commandValue)
			_ = msg.Nack(false, false)
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

		log.Printf("Command handled: %s", envelope.CommandType)
		err = msg.Ack(false) // Acknowledge message
		if err != nil {
			log.Printf("Failed to acknowledge message: %v", err)
		} else {
			log.Printf("Message acknowledged: %s", msg.Body)
		}
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
