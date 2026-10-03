package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"zimbres/uptime-kuma-api/internal/config"
	"zimbres/uptime-kuma-api/internal/models"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

func TestAuthMiddleware_Disabled(t *testing.T) {
	cfg := &config.Config{
		EnableAuth: false,
	}

	reached := false
	r := setupTestRouter()
	r.Use(AuthMiddleware(cfg))
	r.GET("/test", func(c *gin.Context) {
		reached = true
		c.JSON(http.StatusOK, models.APIResponse{Success: true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, reached)
}

func TestAuthMiddleware_Enabled_NoHeader(t *testing.T) {
	cfg := &config.Config{
		EnableAuth: true,
		AuthToken:  "secret-token",
	}

	r := setupTestRouter()
	r.Use(AuthMiddleware(cfg))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.APIResponse{Success: true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Enabled_InvalidFormat(t *testing.T) {
	cfg := &config.Config{
		EnableAuth: true,
		AuthToken:  "secret-token",
	}

	r := setupTestRouter()
	r.Use(AuthMiddleware(cfg))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.APIResponse{Success: true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Enabled_InvalidToken(t *testing.T) {
	cfg := &config.Config{
		EnableAuth: true,
		AuthToken:  "secret-token",
	}

	r := setupTestRouter()
	r.Use(AuthMiddleware(cfg))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.APIResponse{Success: true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Enabled_ValidToken(t *testing.T) {
	cfg := &config.Config{
		EnableAuth: true,
		AuthToken:  "secret-token",
	}

	reached := false
	r := setupTestRouter()
	r.Use(AuthMiddleware(cfg))
	r.GET("/test", func(c *gin.Context) {
		reached = true
		c.JSON(http.StatusOK, models.APIResponse{Success: true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, reached)
}