package domain

import (
	"github.com/jperdior/chatbot-kit/application/event"
	"time"
)

type BaseAggregate struct {
	events    []event.Event
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *BaseAggregate) PullEvents() []event.Event {
	events := a.events
	a.events = []event.Event{}
	return events
}

func (a *BaseAggregate) Record(event event.Event) {
	a.events = append(a.events, event)
}
