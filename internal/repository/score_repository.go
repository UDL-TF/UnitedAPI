package repository

import (
	"fmt"
	"math"

	"github.com/UDL-TF/UnitedAPI/internal/model"
	"gorm.io/gorm"
)

// ScoreRepository handles score-related database operations
type ScoreRepository struct {
	db *gorm.DB
}

// NewScoreRepository creates a new instance of ScoreRepository
func NewScoreRepository(db *gorm.DB) *ScoreRepository {
	return &ScoreRepository{
		db: db,
	}
}

// UpdateMatchRound updates a match round with the score data
func (r *ScoreRepository) UpdateMatchRound(roundID, winnerID, loserID, homeTeamScore, awayTeamScore int) error {
	scoreDifference := math.Abs(float64(homeTeamScore - awayTeamScore))

	result := r.db.Model(&model.MatchRound{}).
		Where("id = ?", roundID).
		Updates(map[string]interface{}{
			"winner_id":        winnerID,
			"loser_id":         loserID,
			"home_team_score":  homeTeamScore,
			"away_team_score":  awayTeamScore,
			"has_outcome":      true,
			"score_difference": scoreDifference,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update match round: %w", result.Error)
	}

	return nil
}

// AreAllRoundsDone checks if all rounds for a match are completed
func (r *ScoreRepository) AreAllRoundsDone(matchID int) (bool, error) {
	// First check if the entire match is marked as not done
	var match model.Match
	result := r.db.Select("manual_not_done").Where("id = ?", matchID).First(&match)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check match status: %w", result.Error)
	}

	if match.ManualNotDone {
		return false, nil
	}

	// Count the number of rounds without outcomes
	var count int64
	result = r.db.Model(&model.MatchRound{}).
		Where("match_id = ? AND has_outcome = ?", matchID, false).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check if all rounds are done: %w", result.Error)
	}

	return count == 0, nil
}

// UpdateMatchStatus updates the status of a match
func (r *ScoreRepository) UpdateMatchStatus(matchID, status int) error {
	result := r.db.Model(&model.Match{}).
		Where("id = ?", matchID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update match status: %w", result.Error)
	}

	return nil
}

// HasTeamReachedRoundWinLimit checks whether either team reached the configured round win limit
func (r *ScoreRepository) HasTeamReachedRoundWinLimit(matchID int) (bool, error) {
	var match model.Match
	matchQuery := r.db.Select("round_win_limit", "home_team_id", "away_team_id").Where("id = ?", matchID).First(&match)
	if matchQuery.Error != nil {
		return false, fmt.Errorf("failed to fetch match info: %w", matchQuery.Error)
	}

	validWinners := make(map[int]struct{})
	if match.RosterHomeID != 0 {
		validWinners[match.RosterHomeID] = struct{}{}
	}
	if match.RosterAwayID != 0 {
		validWinners[match.RosterAwayID] = struct{}{}
	}

	if match.RoundWinLimit <= 0 || len(validWinners) == 0 {
		return false, nil
	}

	var rounds []model.MatchRound
	roundsQuery := r.db.Select("winner_id").Where("match_id = ?", matchID).Find(&rounds)
	if roundsQuery.Error != nil {
		return false, fmt.Errorf("failed to fetch match rounds: %w", roundsQuery.Error)
	}

	wins := make(map[int]int)
	for _, round := range rounds {
		if round.WinnerID == nil {
			continue
		}

		winnerID := *round.WinnerID
		if _, ok := validWinners[winnerID]; !ok {
			continue
		}

		wins[winnerID]++
		if wins[winnerID] >= match.RoundWinLimit {
			return true, nil
		}
	}

	return false, nil
}
