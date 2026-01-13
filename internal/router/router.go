package router

import (
	"github.com/UDL-TF/UnitedAPI/internal/response"
	v1 "github.com/UDL-TF/UnitedAPI/internal/router/v1"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes with versioning
func SetupRouter(engine *gin.Engine) {
	// Health check endpoint (no versioning)
	engine.GET("/health", HealthCheck)

	// API routes
	api := engine.Group("/api")
	{
		v1.RegisterRoutes(api)
	}
}

// HealthCheck handler
func HealthCheck(c *gin.Context) {
	response.OK(c, gin.H{})
}
