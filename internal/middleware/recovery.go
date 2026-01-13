package middleware

import (
	"net/http"

	"github.com/UDL-TF/UnitedAPI/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery middleware recovers from panics and returns 500
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error("Panic recovered", zap.Any("error", err))
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
