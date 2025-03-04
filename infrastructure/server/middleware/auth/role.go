package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

// RoleMiddleware is a middleware that checks if the user has the required roles
func RoleMiddleware(requiredRoles []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, exists := context.Get("claims")
		if !exists {
			context.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			context.Abort()
			return
		}

		roleClaims, ok := claims.(jwt.MapClaims)["roles"].([]interface{})
		if !ok {
			context.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			context.Abort()
			return
		}

		for _, userRole := range roleClaims {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					context.Next()
					return
				}
			}
		}

		context.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
		context.Abort()
	}
}
