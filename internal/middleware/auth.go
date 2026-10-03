package middleware

import (
	"net/http"
	"strings"

	"zimbres/uptime-kuma-api/internal/config"
	"zimbres/uptime-kuma-api/internal/models"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.EnableAuth {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Invalid authorization format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		if parts[1] != cfg.AuthToken {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Invalid token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}