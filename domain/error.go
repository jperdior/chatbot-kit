package domain

import (
	"fmt"
	"log"
)

type DomainError struct {
	Message string
	Key     string
}

// Error implements the error interface for DomainError.
func (e *DomainError) Error() string {
	return fmt.Sprintf("DomainError{Message: %s, Key: %s}", e.Message, e.Key)
}

func NewDomainError(message, key string) *DomainError {
	domainErr := &DomainError{
		Message: message,
		Key:     key,
	}
	log.Println(domainErr) // Log the error
	return domainErr
}

type ErrorResponse struct {
	Error string `json:"error"`
}
