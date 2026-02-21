package handler

import (
	"github.com/UDL-TF/UnitedAPI/internal/logger"
	"github.com/UDL-TF/UnitedAPI/internal/model"
	"github.com/UDL-TF/UnitedAPI/internal/response"
	"github.com/UDL-TF/UnitedAPI/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ScoreHandler handles score-related HTTP requests
type ScoreHandler struct {
	service *service.ScoreService
}

// NewScoreHandler creates a new instance of ScoreHandler
func NewScoreHandler(service *service.ScoreService) *ScoreHandler {
	return &ScoreHandler{
		service: service,
	}
}

// SendScores handles POST /api/v1/scores
// Receives score data from the game server and processes it
func (h *ScoreHandler) SendScores(c *gin.Context) {
	var scoreData model.ScoreData

	// Parse JSON body
	if err := c.ShouldBindJSON(&scoreData); err != nil {
		logger.Log.Warn("Invalid score data received", zap.Error(err))
		response.BadRequest(c, "Error parsing JSON data: "+err.Error())
		return
	}

	// Process the score data
	if err := h.service.ProcessScoreData(&scoreData); err != nil {
		logger.Log.Error("Failed to process score data", zap.Error(err), zap.Any("scoreData", scoreData))
		response.Error(c, 500, "PROCESSING_ERROR", err.Error())
		return
	}

	// Send success response
	response.SuccessWithMessage(c, 200, "Score data received", nil)
}
