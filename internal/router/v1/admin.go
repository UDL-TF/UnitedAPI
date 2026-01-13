package v1

import (
	"github.com/UDL-TF/UnitedAPI/internal/handler"
	"github.com/UDL-TF/UnitedAPI/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes registers admin routes for API v1
func RegisterAdminRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	admin.Use(middleware.Auth())
	admin.Use(middleware.RateLimit(10))
	{
		admin.GET("/stats", handler.GetStats)
		admin.POST("/users/bulk", handler.BulkCreateUsers)
	}
}
