package v1

import (
	"github.com/UDL-TF/UnitedAPI/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterDemoRoutes registers demo-related routes
func RegisterDemoRoutes(rg *gin.RouterGroup) {
	demos := rg.Group("/demos")
	{
		demos.GET("", handler.GetDemo)
	}
}
