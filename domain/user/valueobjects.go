package domain

import (
	"github.com/google/uuid"
)

type UserID struct {
	value uuid.UUID
}

func NewUserID(value uuid.UUID) *UserID {
	return &UserID{value: value}
}

func (u UserID) Value() uuid.UUID {
	return u.value
}

func NewRandomUserID() UserID {
	return UserID{value: uuid.New()}
}
