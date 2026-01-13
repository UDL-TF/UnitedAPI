package middleware

import (
	"strings"

	"github.com/UDL-TF/UnitedAPI/internal/response"
	"github.com/gin-gonic/gin"
)

// Auth middleware for JWT or API key authentication
// This is a placeholder - implement your actual auth logic
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		// Example: Bearer token validation
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// TODO: Implement actual token validation
		// For now, just check if token exists
		if token == "" {
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// Set user info in context after validation
		// c.Set("user_id", userId)

		c.Next()
	}
}
