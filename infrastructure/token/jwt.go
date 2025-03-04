package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTTokenGenerator struct {
	secretKey  string
	expiration int
}

func NewJWTTokenGenerator(secretKey string, expiration int) *JWTTokenGenerator {
	return &JWTTokenGenerator{secretKey: secretKey, expiration: expiration}
}

// TODO this should be in an sso service and work with secrets
func (j *JWTTokenGenerator) GenerateClientToken(clientID, clientName string) (string, error) {
	duration := time.Duration(j.expiration) * 24 * time.Hour
	claims := jwt.MapClaims{
		"client_id":   clientID,
		"client_name": clientName,
		"token_type":  "client",
		"iss":         "user-service",
		"iat":         time.Now().Unix(),
		"exp":         time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
