package v1

import (
	"github.com/UDL-TF/UnitedAPI/internal/config"
	"github.com/UDL-TF/UnitedAPI/internal/handler"
	"github.com/UDL-TF/UnitedAPI/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all v1 routes
func RegisterRoutes(rg *gin.RouterGroup, cfg *config.Config, scoreHandler *handler.ScoreHandler) {
	v1 := rg.Group("/v1")
	v1.Use(middleware.Logger())
	{
		RegisterDemoRoutes(v1, cfg.Auth.SecretDemoPassword)
		RegisterProtectedRoutes(v1)
		RegisterAdminRoutes(v1)
		RegisterScoreRoutes(v1, scoreHandler, cfg.Auth.SecretPassword)
	}
}
