package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/UDL-TF/UnitedAPI/internal/logger"
	"github.com/UDL-TF/UnitedAPI/internal/model"
	"github.com/UDL-TF/UnitedAPI/internal/repository"
	"go.uber.org/zap"
)

// ScoreService handles business logic for score operations
type ScoreService struct {
	repo *repository.ScoreRepository
}

// NewScoreService creates a new instance of ScoreService
func NewScoreService(repo *repository.ScoreRepository) *ScoreService {
	return &ScoreService{
		repo: repo,
	}
}

// ProcessScoreData processes the score data from the game server
func (s *ScoreService) ProcessScoreData(scoreData *model.ScoreData) error {
	logger.Log.Info("Processing score data", zap.Any("scoreData", scoreData))

	// Update the database with match round information
	err := s.repo.UpdateMatchRound(
		scoreData.RoundID,
		scoreData.WinnerTeamID,
		scoreData.LoserTeamID,
		scoreData.HomePoints,
		scoreData.AwayPoints,
	)
	if err != nil {
		logger.Log.Error("Failed to update match round", zap.Error(err))
		return fmt.Errorf("error updating data: %w", err)
	}

	// Call external API to update scores
	if err := s.updateScores(scoreData.MatchID); err != nil {
		if handledErr := s.handleScoreUpdateError(scoreData.MatchID, "Failed to update scores via external API", err); handledErr != nil {
			return handledErr
		}
	}

	// Check if all rounds are done
	allDone, err := s.repo.AreAllRoundsDone(scoreData.MatchID)
	if err != nil {
		logger.Log.Error("Failed to check if all rounds are done", zap.Error(err))
		// Continue processing even if this check fails
	}

	shouldComplete := allDone
	completionReason := "all rounds completed"
	if !shouldComplete {
		reachedLimit, checkErr := s.repo.HasTeamReachedRoundWinLimit(scoreData.MatchID)
		if checkErr != nil {
			logger.Log.Error("Failed to determine if round win limit was reached", zap.Error(checkErr))
		} else if reachedLimit {
			shouldComplete = true
			completionReason = "round win limit reached"
		}
	}

	if shouldComplete {
		logger.Log.Info("Marking match as completed", zap.Int("matchID", scoreData.MatchID), zap.String("reason", completionReason))
		if err := s.completeMatch(scoreData.MatchID); err != nil {
			return err
		}
	}

	logger.Log.Info("Successfully processed score data", zap.Int("matchID", scoreData.MatchID), zap.Int("roundID", scoreData.RoundID))
	return nil
}

// completeMatch finalizes a match by setting its status, determining winner/loser, and notifying the league site
func (s *ScoreService) completeMatch(matchID int) error {
	// Update match status first
	if err := s.repo.UpdateMatchStatus(matchID, 3); err != nil {
		logger.Log.Error("Failed to update match status", zap.Error(err), zap.Int("matchID", matchID))
		return fmt.Errorf("error updating match status: %w", err)
	}

	// Only set winner/loser if ALL rounds are actually completed
	allRoundsDone, err := s.repo.AreAllRoundsDone(matchID)
	if err != nil {
		logger.Log.Error("Failed to check if all rounds are done", zap.Error(err), zap.Int("matchID", matchID))
		// Continue without setting winner/loser if we can't verify all rounds are done
	} else if allRoundsDone {
		// Determine the match winner and loser
		winnerID, loserID, err := s.repo.DetermineMatchWinner(matchID)
		if err != nil {
			logger.Log.Error("Failed to determine match winner", zap.Error(err), zap.Int("matchID", matchID))
		} else {
			// Update match winner and loser information
			if err := s.repo.UpdateMatchWinner(matchID, winnerID, loserID); err != nil {
				logger.Log.Error("Failed to update match winner", zap.Error(err), zap.Int("matchID", matchID), zap.Int("winnerID", winnerID), zap.Int("loserID", loserID))
			} else {
				logger.Log.Info("Match completed with winner and loser", zap.Int("matchID", matchID), zap.Int("winnerID", winnerID), zap.Int("loserID", loserID))
			}
		}
	} else {
		logger.Log.Info("Match completed but not all rounds finished - winner/loser not set", zap.Int("matchID", matchID))
	}

	if err := s.updateScores(matchID); err != nil {
		if handledErr := s.handleScoreUpdateError(matchID, "Failed to update scores via external API after completion", err); handledErr != nil {
			return handledErr
		}
	}

	return nil
}

// updateScores calls the external API to update scores
func (s *ScoreService) updateScores(matchID int) error {
	updateScoresURL := fmt.Sprintf("https://udl.tf/leagues/matches/%d/update_scores", matchID)

	req, err := http.NewRequest(http.MethodPost, updateScoresURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request to update scores: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to update scores: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return &ScoreUpdateError{StatusCode: resp.StatusCode, Body: strings.TrimSpace(string(bodyBytes))}
	}

	return nil
}

// handleScoreUpdateError centralizes how we react to external score sync failures
func (s *ScoreService) handleScoreUpdateError(matchID int, logMessage string, err error) error {
	var scoreErr *ScoreUpdateError
	if errors.As(err, &scoreErr) && scoreErr.StatusCode == http.StatusUnprocessableEntity {
		logger.Log.Warn(logMessage,
			zap.Int("matchID", matchID),
			zap.Int("status", scoreErr.StatusCode),
			zap.String("externalMessage", scoreErr.Body))
		return nil
	}

	logger.Log.Error(logMessage, zap.Error(err), zap.Int("matchID", matchID))
	return fmt.Errorf("error updating scores: %w", err)
}

// ScoreUpdateError captures structured details from the external API when it fails
type ScoreUpdateError struct {
	StatusCode int
	Body       string
}

// Error implements the error interface for ScoreUpdateError
func (e *ScoreUpdateError) Error() string {
	if e == nil {
		return ""
	}
	if e.Body != "" {
		return fmt.Sprintf("external API returned status %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("external API returned status %d", e.StatusCode)
}
