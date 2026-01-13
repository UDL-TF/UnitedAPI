package middleware

import (
	"time"

	"github.com/UDL-TF/UnitedAPI/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger middleware logs request details
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Log.Info("Request processed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
		)
	}
}
