package inmemory

import (
	"context"
	"fmt"
	"github.com/jperdior/chatbot-kit/application/command"
	"log"
)

// CommandBus is an in-memory implementation of the command.Bus.
type CommandBus struct {
	handlers map[command.Type]command.Handler
}

// NewCommandBus initializes a new instance of CommandBus.
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[command.Type]command.Handler),
	}
}

// Dispatch implements the command.Bus interface.
func (b *CommandBus) Dispatch(ctx context.Context, cmd command.Command) error {
	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return nil
	}
	fmt.Print("about to dispatch command\n")
	go func() {
		fmt.Print("Dispatching command\n")
		err := handler.Handle(ctx, cmd)
		if err != nil {
			log.Printf("Error while handling %s - %s\n", cmd.Type(), err)
		}
	}()

	return nil
}

// Register implements the command.Bus interface.
func (b *CommandBus) Register(cmdType command.Type, handler command.Handler) {
	b.handlers[cmdType] = handler
}

// Consume does not apply in in-memory implementation
func (b *CommandBus) Consume() error {
	return nil
}

// Close does not apply in in-memory implementation
func (b *CommandBus) Close() {
	// noop
}
