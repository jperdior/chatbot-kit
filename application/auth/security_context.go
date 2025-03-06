package auth

import (
	"context"
	domain "github.com/jperdior/chatbot-kit/domain/user"
)

type TokenType string

type SecurityContext interface {
	Type() TokenType
	GetIdentifier() string
}

const UserSecurityContextType TokenType = "user"

type UserSecurityContext struct {
	ID    *domain.UserID
	Email string
	Roles []string
}

func (t *UserSecurityContext) GetIdentifier() string {
	return t.ID.String()
}

func (t *UserSecurityContext) Type() TokenType {
	return UserSecurityContextType
}

func NewUserSecurityContext(id *domain.UserID, email string, roles []string) *UserSecurityContext {
	return &UserSecurityContext{
		ID:    id,
		Email: email,
		Roles: roles,
	}
}

const ClientSecurityContextType TokenType = "client"

type ClientSecurityContext struct {
	ClientID   string
	ClientName string
}

func (t *ClientSecurityContext) GetIdentifier() string {
	return t.ClientID
}

func (t *ClientSecurityContext) Type() TokenType {
	return ClientSecurityContextType
}

func NewClientSecurityContext(clientID, clientName string) *ClientSecurityContext {
	return &ClientSecurityContext{
		ClientID:   clientID,
		ClientName: clientName,
	}
}

func GetSecurityContextFromContext(ctx context.Context) SecurityContext {
	securityContext, ok := ctx.Value("securityContext").(SecurityContext)
	if !ok {
		return nil
	}
	return securityContext
}
