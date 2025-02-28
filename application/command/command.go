package command

import (
	"context"
	"encoding/json"
)

// Bus defines the expected behaviour from a command bus.
type Bus interface {
	// Dispatch is the method used to dispatch new commands.
	Dispatch(context.Context, Command) error
	// Register is the method used to register a new command handler.
	Register(Type, Handler)
	// Consume is the method used to listen for new commands.
	Consume() error
	// Close is the method used to close the bus.
	Close()
}

//go:generate mockery --case=snake --outpkg=commandmocks --output=commandmocks --name=Bus

// Type represents an application command type.
type Type string

type CommandEnvelope struct {
	CommandType Type            `json:"type"`
	Data        json.RawMessage `json:"data"`
}

// Command represents an application command.
type Command interface {
	Type() Type
}

// Handler defines the expected behaviour from a command handler.
type Handler interface {
	Handle(context.Context, Command) error
}
