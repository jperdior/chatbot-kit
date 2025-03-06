package auth

import (
	"context"
)

type SecurityProvider interface {
	GetSecurityContext(ctx context.Context) SecurityContext
}

type JWTSecurityProvider struct {
}

func NewJWTSecurityProvider() *JWTSecurityProvider {
	return &JWTSecurityProvider{}
}

func (p *JWTSecurityProvider) GetSecurityContext(ctx context.Context) SecurityContext {
	securityContext, ok := ctx.Value("securityContext").(SecurityContext)
	if !ok {
		return nil
	}
	return securityContext
}
