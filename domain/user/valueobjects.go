package domain

import (
	"context"
	"github.com/jperdior/chatbot-kit/domain"
)

type UserID struct {
	*domain.UUIDValueObject
}

func NewUserID(value string) (*UserID, error) {
	uid, err := domain.NewUuidValueObject(value)
	if err != nil {
		return &UserID{}, NewInvalidUserIDError(value)
	}
	return &UserID{UUIDValueObject: uid}, nil
}

func FromContext(ctx context.Context) (*UserID, error) {
	userID := ctx.Value("ID").(string)
	uid, err := domain.NewUuidValueObject(userID)
	if err != nil {
		return &UserID{}, NewInvalidUserIDError(userID)
	}
	return &UserID{UUIDValueObject: uid}, nil
}

func NewRandomUserID() *UserID {
	uid := domain.NewRandomUUIDValueObject()
	return &UserID{UUIDValueObject: uid}
}
