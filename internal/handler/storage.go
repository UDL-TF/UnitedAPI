package handler

import (
	"fmt"
	"io"
	"strconv"
	"time"

	appctx "github.com/UDL-TF/UnitedAPI/internal/context"
	"github.com/UDL-TF/UnitedAPI/internal/response"
	"github.com/gin-gonic/gin"
)

// UploadFile handles file upload to MinIO
func UploadFile(c *gin.Context) {
	// Get app context from middleware
	ctx := c.MustGet("app_context").(*appctx.AppContext)
	minioClient := ctx.GetMinIOClient()

	if minioClient == nil {
		response.InternalServerError(c, "Storage service not available")
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Failed to get file from form")
		return
	}
	defer file.Close()

	// Get optional object name from form, default to original filename
	objectName := c.PostForm("object_name")
	if objectName == "" {
		objectName = header.Filename
	}

	// Get file size
	fileSize := header.Size

	// Determine content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload file to MinIO
	err = minioClient.UploadFile(objectName, file, fileSize, contentType)
	if err != nil {
		response.InternalServerError(c, "Failed to upload file")
		return
	}

	response.OKWithMessage(c, "File uploaded successfully", gin.H{
		"object_name":  objectName,
		"size":         fileSize,
		"content_type": contentType,
	})
}

// DownloadFile handles file download from MinIO
func DownloadFile(c *gin.Context) {
	// Get app context from middleware
	ctx := c.MustGet("app_context").(*appctx.AppContext)
	minioClient := ctx.GetMinIOClient()

	if minioClient == nil {
		response.InternalServerError(c, "Storage service not available")
		return
	}

	// Get object name from URL parameter
	objectName := c.Param("objectName")
	if objectName == "" {
		response.BadRequest(c, "Object name is required")
		return
	}

	// Get file info first to check if file exists and get metadata
	fileInfo, err := minioClient.GetFileInfo(objectName)
	if err != nil {
		response.NotFound(c, "File not found")
		return
	}

	// Download file from MinIO
	object, err := minioClient.DownloadFile(objectName)
	if err != nil {
		response.InternalServerError(c, "Failed to download file")
		return
	}
	defer object.Close()

	// Set response headers
	c.Header("Content-Type", fileInfo.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size, 10))

	// Copy file content to response
	_, err = io.Copy(c.Writer, object)
	if err != nil {
		response.InternalServerError(c, "Failed to stream file")
		return
	}
}

// DeleteFile handles file deletion from MinIO
func DeleteFile(c *gin.Context) {
	// Get app context from middleware
	ctx := c.MustGet("app_context").(*appctx.AppContext)
	minioClient := ctx.GetMinIOClient()

	if minioClient == nil {
		response.InternalServerError(c, "Storage service not available")
		return
	}

	// Get object name from URL parameter
	objectName := c.Param("objectName")
	if objectName == "" {
		response.BadRequest(c, "Object name is required")
		return
	}

	// Delete file from MinIO
	err := minioClient.DeleteFile(objectName)
	if err != nil {
		response.InternalServerError(c, "Failed to delete file")
		return
	}

	response.OKWithMessage(c, "File deleted successfully", gin.H{
		"object_name": objectName,
	})
}

// ListFiles handles listing files in MinIO
func ListFiles(c *gin.Context) {
	// Get app context from middleware
	ctx := c.MustGet("app_context").(*appctx.AppContext)
	minioClient := ctx.GetMinIOClient()

	if minioClient == nil {
		response.InternalServerError(c, "Storage service not available")
		return
	}

	// Get optional prefix parameter
	prefix := c.Query("prefix")

	// List files from MinIO
	objects, err := minioClient.ListFiles(prefix)
	if err != nil {
		response.InternalServerError(c, "Failed to list files")
		return
	}

	// Format response
	files := make([]gin.H, len(objects))
	for i, obj := range objects {
		files[i] = gin.H{
			"name":          obj.Key,
			"size":          obj.Size,
			"content_type":  obj.ContentType,
			"last_modified": obj.LastModified,
			"etag":          obj.ETag,
		}
	}

	response.OKWithMessage(c, "Files listed successfully", gin.H{
		"files": files,
		"count": len(files),
	})
}

// GetPresignedURL generates a presigned URL for direct access
func GetPresignedURL(c *gin.Context) {
	// Get app context from middleware
	ctx := c.MustGet("app_context").(*appctx.AppContext)
	minioClient := ctx.GetMinIOClient()

	if minioClient == nil {
		response.InternalServerError(c, "Storage service not available")
		return
	}

	// Get object name from URL parameter
	objectName := c.Param("objectName")
	if objectName == "" {
		response.BadRequest(c, "Object name is required")
		return
	}

	// Get method from query parameter (default to GET)
	method := c.DefaultQuery("method", "GET")
	if method != "GET" && method != "PUT" {
		response.BadRequest(c, "Method must be GET or PUT")
		return
	}

	// Get expiry duration from query parameter (default to 1 hour)
	expiryStr := c.DefaultQuery("expiry", "1h")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		response.BadRequest(c, "Invalid expiry duration")
		return
	}

	// Generate presigned URL
	url, err := minioClient.GetPresignedURL(objectName, method, expiry)
	if err != nil {
		response.InternalServerError(c, "Failed to generate presigned URL")
		return
	}

	response.OKWithMessage(c, "Presigned URL generated successfully", gin.H{
		"url":         url,
		"method":      method,
		"expires":     time.Now().Add(expiry),
		"object_name": objectName,
	})
}
