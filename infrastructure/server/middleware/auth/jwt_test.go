package auth

import (
	"chatbot-service/internal/platform/server/handler/status"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func generateValidToken(secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   "123",
		"roles": []string{"ROLE_USER", "ROLE_SUPER_ADMIN"},
		"exp":   time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
		"iat":   time.Now().Unix(),                    // Issued at time
		"iss":   "user-service",                       // Issuer
	})
	return token.SignedString([]byte(secretKey))
}

func TestJWTMiddleware(t *testing.T) {

	secretKey := "secret"
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(JWTMiddleware(secretKey))

	engine.GET("/status", status.StatusHandler())

	t.Run("when the Authorization header is missing", func(t *testing.T) {

		// Setting up the HTTP recorder and the request
		httpRecorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/status", nil)

		// Performing the request
		engine.ServeHTTP(httpRecorder, req)

		// Asserting the response status code
		assert.Equal(t, http.StatusUnauthorized, httpRecorder.Code)
	})

	t.Run("when the Authorization header is present but is invalid", func(t *testing.T) {

		// Setting up the HTTP recorder and the request
		httpRecorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/status", nil)
		req.Header.Add("Authorization", "Bearer token")

		// Performing the request
		engine.ServeHTTP(httpRecorder, req)

		// Asserting the response status code
		assert.Equal(t, http.StatusUnauthorized, httpRecorder.Code)
	})

	t.Run("when the Authorization header is present and is valid", func(t *testing.T) {

		validToken, err := generateValidToken(secretKey)
		assert.NoError(t, err)

		// Setting up the HTTP recorder and the request
		httpRecorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/status", nil)
		req.Header.Add("Authorization", "Bearer "+validToken)

		// Performing the request
		engine.ServeHTTP(httpRecorder, req)

		// Asserting the response status code
		assert.Equal(t, http.StatusOK, httpRecorder.Code)
	})
}
