package auth

import (
	"chatbot-service/internal/platform/server/handler/status"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setClaimsMiddleware(claims jwt.MapClaims) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	}
}

func TestRoleMiddleware(t *testing.T) {

	t.Run("Token with required roles", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		engine := gin.New()

		claims := jwt.MapClaims{
			"roles": []interface{}{"ROLE_SUPER_ADMIN"},
		}
		engine.Use(setClaimsMiddleware(claims))
		engine.Use(RoleMiddleware([]string{"ROLE_SUPER_ADMIN"}))
		engine.GET("/status", status.StatusHandler())

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/status", nil)
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("Token without required roles", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		engine := gin.New()

		claims := jwt.MapClaims{
			"roles": []interface{}{"ROLE_USER"},
		}
		engine.Use(setClaimsMiddleware(claims))
		engine.Use(RoleMiddleware([]string{"ROLE_SUPER_ADMIN"}))
		engine.GET("/status", status.StatusHandler())

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/status", nil)
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})
}
