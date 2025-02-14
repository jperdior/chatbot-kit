package domain

import (
	"github.com/jperdior/chatbot-kit/domain"
)

type UserID struct {
	domain.UUIDValueObject
}

func NewUserID(value string) (UserID, error) {
	uid, err := domain.NewUuidValueObject(value)
	if err != nil {
		return UserID{}, err
	}
	return UserID{UUIDValueObject: *uid}, nil
}

func NewRandomUserID() UserID {
	uid := domain.NewRandomUUIDValueObject()
	return UserID{UUIDValueObject: *uid}
}
