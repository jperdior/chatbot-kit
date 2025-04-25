package user

import (
	"github.com/jperdior/chatbot-kit/application/event"
)

const UserRegisteredType event.Type = "user.user_registered"

type UserRegisteredEvent struct {
	*event.BaseEvent `json:",inline"`
	Email            string   `json:"email"`
	Name             string   `json:"name"`
	Roles            []string `json:"roles"`
}

func NewUserRegisteredEvent(id, email, name string, roles []string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: event.NewBaseEvent(id),
		Email:     email,
		Name:      name,
		Roles:     roles,
	}
}

func (e *UserRegisteredEvent) Type() event.Type {
	return UserRegisteredType
}

const UserDeletedType event.Type = "user.user_deleted"

type UserDeletedEvent struct {
	*event.BaseEvent
}

func NewUserDeletedEvent(id string) *UserDeletedEvent {
	return &UserDeletedEvent{
		BaseEvent: event.NewBaseEvent(id),
	}
}

func (e *UserDeletedEvent) Type() event.Type {
	return UserDeletedType
}
