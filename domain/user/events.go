package domain

import (
	"github.com/jperdior/chatbot-kit/application/event"
)

const UserRegisteredType event.Type = "user.user_registered"

type UserRegisteredEvent struct {
	*event.BaseEvent
	email string
	roles []string
}

func NewUserRegisteredEvent(id string, email string, roles []string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: event.NewBaseEvent(id),
		email:     email,
		roles:     roles,
	}
}

func (e *UserRegisteredEvent) Email() string {
	return e.email
}

func (e *UserRegisteredEvent) Roles() []string {
	return e.roles
}

func (e *UserRegisteredEvent) Type() event.Type {
	return UserRegisteredType
}
