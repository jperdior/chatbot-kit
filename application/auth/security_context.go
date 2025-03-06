package auth

import (
	domain "github.com/jperdior/chatbot-kit/domain/user"
)

type SecurityContextType string

type SecurityContext interface {
	Type() SecurityContextType
	GetIdentifier() string
}

const UserSecurityContextType SecurityContextType = "user"

type UserSecurityContext struct {
	ID    *domain.UserID
	Email string
	Roles []string
}

func (t *UserSecurityContext) GetIdentifier() string {
	return t.ID.String()
}

func (t *UserSecurityContext) Type() SecurityContextType {
	return UserSecurityContextType
}

func NewUserSecurityContext(id *domain.UserID, email string, roles []string) *UserSecurityContext {
	return &UserSecurityContext{
		ID:    id,
		Email: email,
		Roles: roles,
	}
}

const ClientSecurityContextType SecurityContextType = "client"

type ClientSecurityContext struct {
	ClientID   string
	ClientName string
}

func (t *ClientSecurityContext) GetIdentifier() string {
	return t.ClientID
}

func (t *ClientSecurityContext) Type() SecurityContextType {
	return ClientSecurityContextType
}

func NewClientSecurityContext(clientID, clientName string) *ClientSecurityContext {
	return &ClientSecurityContext{
		ClientID:   clientID,
		ClientName: clientName,
	}
}
