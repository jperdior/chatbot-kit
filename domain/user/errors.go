package domain

import "github.com/jperdior/chatbot-kit/domain"

type InvalidUserIDError struct {
	*domain.DomainError
}

func NewInvalidUserIDError(userID string) *InvalidUserIDError {
	return &InvalidUserIDError{
		DomainError: &domain.DomainError{
			Message: "User with ID " + userID + " is not valid",
			Key:     "invalid_user_id",
		},
	}
}

func (e InvalidUserIDError) Error() string {
	return e.Message
}
