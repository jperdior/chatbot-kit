package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// Bus defines the expected behaviour from an event bus.
type Bus interface {
	// Publish is the method used to publish new events.
	Publish([]Event) error
	//Subscribe is the method used to subscribe to an event.
	Subscribe(Type, Handler)
}

//go:generate mockery --case=snake --outpkg=eventmocks --output=eventmocks --name=Bus

// Handler defines the expected behaviour from an event handler.
type Handler interface {
	Handle(Event) error
}

// Type represents a domain event type.
type Type string

// EventDTO represents a data transfer object for an event.
type EventDTO interface{}

// EventEnvelope represents an envelope for an event.
type EventEnvelope struct {
	EventType Type            `json:"type"` // The type of the event
	Data      json.RawMessage `json:"data"`
}

// Event represents a domain event.
type Event interface {
	ID() string
	AggregateID() string
	OccurredOn() time.Time
	Type() Type
}

type BaseEvent struct {
	eventID     string
	aggregateID string
	occurredOn  time.Time
}

func NewBaseEvent(aggregateID string) BaseEvent {
	return BaseEvent{
		eventID:     uuid.New().String(),
		aggregateID: aggregateID,
		occurredOn:  time.Now(),
	}
}

func (b BaseEvent) ID() string {
	return b.eventID
}

func (b BaseEvent) OccurredOn() time.Time {
	return b.occurredOn
}

func (b BaseEvent) AggregateID() string {
	return b.aggregateID
}
