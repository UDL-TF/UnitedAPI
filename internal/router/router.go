package router

import (
	"github.com/UDL-TF/UnitedAPI/internal/config"
	"github.com/UDL-TF/UnitedAPI/internal/handler"
	"github.com/UDL-TF/UnitedAPI/internal/repository"
	"github.com/UDL-TF/UnitedAPI/internal/response"
	v1 "github.com/UDL-TF/UnitedAPI/internal/router/v1"
	"github.com/UDL-TF/UnitedAPI/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter configures all routes with versioning
func SetupRouter(engine *gin.Engine, cfg *config.Config, db *gorm.DB) {
	// Health check endpoint (no versioning)
	engine.GET("/health", HealthCheck)

	scoreRepo := repository.NewScoreRepository(db)
	scoreService := service.NewScoreService(scoreRepo)
	scoreHandler := handler.NewScoreHandler(scoreService)

	// API routes
	api := engine.Group("/api")
	{
		v1.RegisterRoutes(api, cfg, scoreHandler)
	}
}

// HealthCheck handler
func HealthCheck(c *gin.Context) {
	response.OK(c, gin.H{})
}
