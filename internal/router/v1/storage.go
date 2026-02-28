package v1

import (
	"github.com/UDL-TF/UnitedAPI/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterStorageRoutes registers storage-related routes
func RegisterStorageRoutes(rg *gin.RouterGroup) {
	storage := rg.Group("/storage")
	{
		// File upload
		storage.POST("/upload", handler.UploadFile)

		// File download
		storage.GET("/download/:objectName", handler.DownloadFile)

		// File deletion
		storage.DELETE("/delete/:objectName", handler.DeleteFile)

		// List files
		storage.GET("/files", handler.ListFiles)

		// Generate presigned URL
		storage.GET("/presigned/:objectName", handler.GetPresignedURL)
	}
}
