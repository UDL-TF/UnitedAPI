package model

// ScoreData represents the score data sent from the game server
type ScoreData struct {
	MatchID      int `json:"match_id" binding:"required"`
	RoundID      int `json:"round_id" binding:"required"`
	WinnerTeamID int `json:"winner_team_id" binding:"required"`
	LoserTeamID  int `json:"loser_team_id" binding:"required"`
	AwayPoints   int `json:"away_points" binding:"required"`
	HomePoints   int `json:"home_points" binding:"required"`
}

// MatchRound represents a match round in the database
type MatchRound struct {
	ID              int     `gorm:"primaryKey"`
	MatchID         int     `gorm:"column:match_id"`
	MapID           int     `gorm:"column:map_id"`
	HomeTeamScore   int     `gorm:"column:home_team_score"`
	AwayTeamScore   int     `gorm:"column:away_team_score"`
	LoserID         *int    `gorm:"column:loser_id"`
	WinnerID        *int    `gorm:"column:winner_id"`
	HasOutcome      bool    `gorm:"column:has_outcome"`
	ScoreDifference float32 `gorm:"column:score_difference"`
	HomeReady       bool    `gorm:"column:home_ready"`
	AwayReady       bool    `gorm:"column:away_ready"`
}

// TableName overrides the table name for MatchRound
func (MatchRound) TableName() string {
	return "league_match_rounds"
}

// Match represents a match in the database
type Match struct {
	ID            int  `gorm:"primaryKey"`
	Status        int  `gorm:"column:status"`
	ManualNotDone bool `gorm:"column:manual_not_done"`
	RosterAwayID  int  `gorm:"column:away_team_id"`
	RosterHomeID  int  `gorm:"column:home_team_id"`
	WinLimit      int  `gorm:"column:win_limit"`
	RoundWinLimit int  `gorm:"column:round_win_limit"`
}

// TableName overrides the table name for Match
func (Match) TableName() string {
	return "league_matches"
}
