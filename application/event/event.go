package event

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// Bus defines the expected behaviour from an event bus.
type Bus interface {
	// Publish is the method used to publish new events.
	Publish(context.Context, []Event) error
	//Subscribe is the method used to subscribe to an event.
	Subscribe(Type, Handler)
	BindQueue(queue, routingKey string) error
	Consume(queue string) error
	Close()
}

//go:generate mockery --case=snake --outpkg=eventmocks --output=eventmocks --name=Bus

// Handler defines the expected behaviour from an event handler.
type Handler interface {
	Handle(context.Context, Event) error
}

// Type represents a domain event type.
type Type string

// EventEnvelope represents an envelope for an event.
type EventEnvelope struct {
	EventType Type            `json:"type"` // The type of the event
	Data      json.RawMessage `json:"data"`
}

// Event represents a domain event.
type Event interface {
	ID() string
	GetAggregateID() string
	GetOccurredOn() time.Time
	Type() Type
}

type BaseEvent struct {
	EventID     string    `json:"id"`
	AggregateID string    `json:"aggregate_id"`
	OccurredOn  time.Time `json:"occurred_on"`
}

func NewBaseEvent(aggregateID string) *BaseEvent {
	return &BaseEvent{
		EventID:     uuid.New().String(),
		AggregateID: aggregateID,
		OccurredOn:  time.Now(),
	}
}

func (b *BaseEvent) ID() string {
	return b.EventID
}

func (b *BaseEvent) GetOccurredOn() time.Time {
	return b.OccurredOn
}

func (b *BaseEvent) GetAggregateID() string {
	return b.AggregateID
}
