package v1

import (
	"github.com/UDL-TF/UnitedAPI/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers user routes for API v1
func RegisterUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("", handler.GetUsers)
		users.GET("/:id", handler.GetUserByID)
		users.POST("", handler.CreateUser)
		users.PUT("/:id", handler.UpdateUser)
		users.DELETE("/:id", handler.DeleteUser)
	}
}
