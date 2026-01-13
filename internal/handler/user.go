package handler

import (
	"github.com/UDL-TF/UnitedAPI/internal/middleware"
	"github.com/UDL-TF/UnitedAPI/internal/response"
	"github.com/gin-gonic/gin"
)

// GetUsers handles GET /api/v1/users
func GetUsers(c *gin.Context) {
	// Get app context
	appCtx := middleware.GetAppContext(c)

	// Example: You can access DB from context
	// db := appCtx.GetDB()
	// users, err := userService.FindAll(db)

	// For now, return mock data
	_ = appCtx // Placeholder to avoid unused variable error

	users := []gin.H{
		{"id": 1, "name": "John Doe"},
		{"id": 2, "name": "Jane Smith"},
	}

	response.OK(c, gin.H{"users": users})
}

// GetUserByID handles GET /api/v1/users/:id
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	appCtx := middleware.GetAppContext(c)
	_ = appCtx // Placeholder

	user := gin.H{
		"id":    id,
		"name":  "John Doe",
		"email": "john@example.com",
	}

	response.OK(c, user)
}

// CreateUser handles POST /api/v1/users
func CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	appCtx := middleware.GetAppContext(c)
	_ = appCtx // Placeholder

	// Example: Create user using service
	// db := appCtx.GetDB()
	// user, err := userService.Create(db, req)

	user := gin.H{
		"id":    1,
		"name":  req.Name,
		"email": req.Email,
	}

	response.Created(c, user)
}

// UpdateUser handles PUT /api/v1/users/:id
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

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

	user := gin.H{
		"id":      id,
		"name":    req.Name,
		"email":   req.Email,
		"updated": true,
	}

	response.OKWithMessage(c, "User updated successfully", user)
}

// DeleteUser handles DELETE /api/v1/users/:id
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	appCtx := middleware.GetAppContext(c)
	_ = appCtx // Placeholder

	result := gin.H{
		"id":      id,
		"deleted": true,
	}

	response.OKWithMessage(c, "User deleted successfully", result)
}
