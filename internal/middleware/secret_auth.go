package middleware

import (
	"github.com/UDL-TF/UnitedAPI/internal/response"
	"github.com/gin-gonic/gin"
)

// SecretPasswordAuth middleware checks for secret_password query parameter
func SecretPasswordAuth(secretPassword string) gin.HandlerFunc {
	return func(c *gin.Context) {
		querySecretPassword := c.Query("secret_password")

		if querySecretPassword == "" {
			response.Unauthorized(c, "Secret password is required")
			c.Abort()
			return
		}

		if querySecretPassword != secretPassword {
			response.Unauthorized(c, "Invalid secret password")
			c.Abort()
			return
		}

		c.Next()
	}
}
