package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jperdior/chatbot-kit/application/auth"
	domain "github.com/jperdior/chatbot-kit/domain/user"
	"net/http"
)

// JWTMiddleware is a middleware that checks for a valid JWT token in the Authorization header
func JWTMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Identify token type
		tokenType, _ := claims["token_type"].(string)
		c.Set("tokenType", tokenType)
		// If it's a user token, extract user-specific claims
		if tokenType == "user" {
			userID, err := domain.NewUserID(claims["ID"].(string))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
				c.Abort()
				return
			}
			roles, ok := claims["roles"].([]string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user roles"})
				c.Abort()
				return
			}
			userSecurityContext := auth.NewUserSecurityContext(
				userID,
				claims["email"].(string),
				roles,
			)
			c.Set("securityContext", userSecurityContext)
		} else if tokenType == "client" {
			clientSecurityContext := auth.NewClientSecurityContext(
				claims["client_id"].(string),
				claims["client_name"].(string),
			)
			c.Set("securityContext", clientSecurityContext)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unknown token type"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("authToken", tokenString)
		c.Next()
	}
}
