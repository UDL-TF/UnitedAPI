package v1

import (
	"github.com/UDL-TF/UnitedAPI/internal/handler"
	"github.com/UDL-TF/UnitedAPI/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterScoreRoutes registers score-related routes
func RegisterScoreRoutes(rg *gin.RouterGroup, scoreHandler *handler.ScoreHandler, secretPassword string) {
	scores := rg.Group("/scores")
	{
		// POST /api/v1/scores - requires secret_password query parameter
		scores.POST("", middleware.SecretPasswordAuth(secretPassword), scoreHandler.SendScores)
	}
}
