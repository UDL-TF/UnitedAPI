package v1

import (
	"github.com/UDL-TF/UnitedAPI/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all v1 routes
func RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.Use(middleware.Logger())
	{
		RegisterUserRoutes(v1)
		RegisterProtectedRoutes(v1)
		RegisterAdminRoutes(v1)
	}
}
