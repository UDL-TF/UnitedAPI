package v1

import (
	"github.com/UDL-TF/UnitedAPI/internal/handler"
	"github.com/UDL-TF/UnitedAPI/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterDemoRoutes registers demo-related routes
func RegisterDemoRoutes(rg *gin.RouterGroup, uploadSecretPassword string) {
	demos := rg.Group("/demos")
	{
		// Public routes - for getting demo information
		demos.GET("", handler.GetDemo)                   // Get demo info by query params
		demos.HEAD("", handler.GetDemo)                  // Support HEAD requests for demo availability checks
		demos.GET("/list", handler.ListDemos)            // List demos with pagination
		demos.GET("/download/:id", handler.DownloadDemo) // Download demo file by ID

		// Upload route - requires secret_password query parameter
		demos.POST("/upload", middleware.SecretPasswordAuth(uploadSecretPassword), handler.UploadDemo)
	}
}
