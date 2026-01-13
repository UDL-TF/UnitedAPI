package v1

import (
	"github.com/UDL-TF/UnitedAPI/internal/handler"
	"github.com/UDL-TF/UnitedAPI/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterProtectedRoutes registers protected routes for API v1
func RegisterProtectedRoutes(rg *gin.RouterGroup) {
	protected := rg.Group("/protected")
	protected.Use(middleware.Auth())
	{
		protected.GET("/profile", handler.GetProfile)
		protected.PUT("/profile", handler.UpdateProfile)
	}
}
