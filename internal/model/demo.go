package model

import (
	"fmt"
	"time"
)

// Demo represents a demo file in the database
type Demo struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	IsTournament   bool      `gorm:"column:is_tournament" json:"is_tournament"`
	TournamentID   *int      `gorm:"column:tournament_id" json:"tournament_id,omitempty"`
	MatchID        *int      `gorm:"column:match_id" json:"match_id,omitempty"`
	RoundID        *int      `gorm:"column:round_id" json:"round_id,omitempty"`
	RawDemoName    string    `gorm:"column:raw_demo_name;not null" json:"raw_demo_name"`
	FileSize       int64     `gorm:"column:file_size;not null" json:"file_size"`
	BlueTeamName   string    `gorm:"column:blue_team_name;not null" json:"blue_team_name"`
	RedTeamName    string    `gorm:"column:red_team_name;not null" json:"red_team_name"`
	ObjectName     string    `gorm:"column:object_name;not null" json:"object_name"`
	ContentType    string    `gorm:"column:content_type" json:"content_type"`
	IsCompressed   bool      `gorm:"column:is_compressed;default:true" json:"is_compressed"`
	UploadedAt     time.Time `gorm:"column:uploaded_at;autoCreateTime" json:"uploaded_at"`
	CompressedSize *int64    `gorm:"column:compressed_size" json:"compressed_size,omitempty"`
}

// TableName overrides the table name for Demo
func (Demo) TableName() string {
	return "demos"
}

// DemoUploadRequest represents the request payload for uploading a demo
type DemoUploadRequest struct {
	IsTournament bool   `form:"is_tournament" binding:"required"`
	TournamentID *int   `form:"tournament_id"`
	MatchID      *int   `form:"match_id"`
	RoundID      *int   `form:"round_id"`
	RawDemoName  string `form:"raw_demo_name" binding:"required"`
	BlueTeamName string `form:"blue_team_name" binding:"required"`
	RedTeamName  string `form:"red_team_name" binding:"required"`
}

// DemoGetRequest represents the query parameters for getting a demo
type DemoGetRequest struct {
	MatchID      *int `form:"match_id"`
	RoundID      *int `form:"round_id"`
	TournamentID *int `form:"tournament_id"`
}

// ValidateTournamentData validates tournament-specific fields
func (d *DemoUploadRequest) ValidateTournamentData() error {
	if d.IsTournament {
		if d.TournamentID == nil || d.MatchID == nil || d.RoundID == nil {
			return fmt.Errorf("tournament_id, match_id, and round_id are required for tournament demos")
		}
	}
	return nil
}
