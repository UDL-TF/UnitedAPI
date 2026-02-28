package handler

import (
	"fmt"
	"io"
	"strconv"

	appctx "github.com/UDL-TF/UnitedAPI/internal/context"
	"github.com/UDL-TF/UnitedAPI/internal/model"
	"github.com/UDL-TF/UnitedAPI/internal/repository"
	"github.com/UDL-TF/UnitedAPI/internal/response"
	"github.com/UDL-TF/UnitedAPI/internal/service"
	"github.com/gin-gonic/gin"
)

// UploadDemo handles demo file upload with compression and metadata
func UploadDemo(c *gin.Context) {
	// Get app context from middleware
	ctx := c.MustGet("app_context").(*appctx.AppContext)
	db := ctx.GetDB()
	minioClient := ctx.GetMinIOClient()

	if db == nil {
		response.InternalServerError(c, "Database not available")
		return
	}

	if minioClient == nil {
		response.InternalServerError(c, "Storage service not available")
		return
	}

	// Parse form data
	var req model.DemoUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, fmt.Sprintf("Invalid request data: %v", err))
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("demo_file")
	if err != nil {
		response.BadRequest(c, "Demo file is required")
		return
	}
	defer file.Close()

	// Create service
	demoRepo := repository.NewDemoRepository(db)
	demoService := service.NewDemoService(demoRepo, minioClient)

	// Upload demo
	demo, err := demoService.UploadDemo(&req, file, header)
	if err != nil {
		response.InternalServerError(c, fmt.Sprintf("Failed to upload demo: %v", err))
		return
	}

	response.OKWithMessage(c, "Demo uploaded successfully", demo)
}

// GetDemo handles demo download and info retrieval
func GetDemo(c *gin.Context) {
	// Get app context from middleware
	ctx := c.MustGet("app_context").(*appctx.AppContext)
	db := ctx.GetDB()

	if db == nil {
		response.InternalServerError(c, "Database not available")
		return
	}

	// Parse query parameters
	var req model.DemoGetRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, fmt.Sprintf("Invalid query parameters: %v", err))
		return
	}

	// Create service
	demoRepo := repository.NewDemoRepository(db)
	demoService := service.NewDemoService(demoRepo, ctx.GetMinIOClient())

	// Get demo information
	demo, err := demoService.GetDemo(&req)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OKWithMessage(c, "Demo found", demo)
}

// DownloadDemo handles demo file download
func DownloadDemo(c *gin.Context) {
	// Get app context from middleware
	ctx := c.MustGet("app_context").(*appctx.AppContext)
	db := ctx.GetDB()
	minioClient := ctx.GetMinIOClient()

	if db == nil || minioClient == nil {
		response.InternalServerError(c, "Required services not available")
		return
	}

	// Get demo ID from URL parameter
	demoIDStr := c.Param("id")
	demoID, err := strconv.Atoi(demoIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid demo ID")
		return
	}

	// Create service
	demoRepo := repository.NewDemoRepository(db)
	demoService := service.NewDemoService(demoRepo, minioClient)

	// Download demo
	object, demo, err := demoService.DownloadDemo(demoID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	defer object.Close()

	// Set response headers for compressed file download
	filename := fmt.Sprintf("%s.zst", demo.RawDemoName)
	c.Header("Content-Type", demo.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	if demo.CompressedSize != nil {
		c.Header("Content-Length", fmt.Sprintf("%d", *demo.CompressedSize))
	}

	// Copy file content to response
	_, err = io.Copy(c.Writer, object)
	if err != nil {
		response.InternalServerError(c, "Failed to stream file")
		return
	}
}

// ListDemos handles demo listing with pagination
func ListDemos(c *gin.Context) {
	// Get app context from middleware
	ctx := c.MustGet("app_context").(*appctx.AppContext)
	db := ctx.GetDB()

	if db == nil {
		response.InternalServerError(c, "Database not available")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// Create service
	demoRepo := repository.NewDemoRepository(db)
	demoService := service.NewDemoService(demoRepo, ctx.GetMinIOClient())

	// List demos
	demos, err := demoService.ListDemos(page, limit)
	if err != nil {
		response.InternalServerError(c, fmt.Sprintf("Failed to list demos: %v", err))
		return
	}

	response.OKWithMessage(c, "Demos retrieved successfully", gin.H{
		"demos": demos,
		"page":  page,
		"limit": limit,
		"count": len(demos),
	})
}
