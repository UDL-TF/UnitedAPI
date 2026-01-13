package handler

import (
	"github.com/UDL-TF/UnitedAPI/internal/middleware"
	"github.com/UDL-TF/UnitedAPI/internal/response"
	"github.com/gin-gonic/gin"
)

// GetProfile handles GET /api/v1/protected/profile
func GetProfile(c *gin.Context) {
	appCtx := middleware.GetAppContext(c)
	_ = appCtx // Placeholder

	// In a real app, you'd get the user ID from the JWT token
	// userId := c.GetString("user_id")

	profile := gin.H{
		"id":    1,
		"name":  "John Doe",
		"email": "john@example.com",
		"role":  "user",
	}

	response.OK(c, profile)
}

// UpdateProfile handles PUT /api/v1/protected/profile
func UpdateProfile(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	appCtx := middleware.GetAppContext(c)
	_ = appCtx // Placeholder

	profile := gin.H{
		"name":    req.Name,
		"email":   req.Email,
		"updated": true,
	}

	response.OKWithMessage(c, "Profile updated successfully", profile)
}

// GetStats handles GET /api/v1/admin/stats
func GetStats(c *gin.Context) {
	appCtx := middleware.GetAppContext(c)
	_ = appCtx // Placeholder

	stats := gin.H{
		"total_users":    100,
		"active_users":   75,
		"total_requests": 1000,
	}

	response.OK(c, stats)
}

// BulkCreateUsers handles POST /api/v1/admin/users/bulk
func BulkCreateUsers(c *gin.Context) {
	var req struct {
		Users []struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		} `json:"users" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	appCtx := middleware.GetAppContext(c)
	_ = appCtx // Placeholder

	result := gin.H{
		"created": len(req.Users),
	}

	response.OKWithMessage(c, "Users created successfully", result)
}
