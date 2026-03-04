package model

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestScoreDataValidation(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		payload  map[string]interface{}
		wantErr  bool
		errField string
	}{
		{
			name: "Valid score data with zero away_points",
			payload: map[string]interface{}{
				"match_id":       497,
				"round_id":       494,
				"winner_team_id": 204,
				"loser_team_id":  268,
				"home_points":    10,
				"away_points":    0, // Should be valid now
			},
			wantErr: false,
		},
		{
			name: "Valid score data with positive values",
			payload: map[string]interface{}{
				"match_id":       497,
				"round_id":       494,
				"winner_team_id": 204,
				"loser_team_id":  268,
				"home_points":    16,
				"away_points":    12,
			},
			wantErr: false,
		},
		{
			name: "Invalid - negative away_points",
			payload: map[string]interface{}{
				"match_id":       497,
				"round_id":       494,
				"winner_team_id": 204,
				"loser_team_id":  268,
				"home_points":    10,
				"away_points":    -1, // Should fail validation
			},
			wantErr:  true,
			errField: "AwayPoints",
		},
		{
			name: "Invalid - negative home_points",
			payload: map[string]interface{}{
				"match_id":       497,
				"round_id":       494,
				"winner_team_id": 204,
				"loser_team_id":  268,
				"home_points":    -5, // Should fail validation
				"away_points":    12,
			},
			wantErr:  true,
			errField: "HomePoints",
		},
		{
			name: "Invalid - missing required field match_id",
			payload: map[string]interface{}{
				"round_id":       494,
				"winner_team_id": 204,
				"loser_team_id":  268,
				"home_points":    10,
				"away_points":    0,
			},
			wantErr:  true,
			errField: "MatchID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test HTTP request
			jsonBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create a test response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test binding
			var scoreData ScoreData
			err := c.ShouldBindJSON(&scoreData)

			if tt.wantErr {
				assert.Error(t, err, "Expected validation error for test case: %s", tt.name)
				if tt.errField != "" {
					assert.Contains(t, err.Error(), tt.errField, "Error should mention field %s", tt.errField)
				}
			} else {
				assert.NoError(t, err, "Expected no validation error for test case: %s", tt.name)
				// Verify the data was bound correctly
				assert.Equal(t, tt.payload["match_id"], scoreData.MatchID)
				assert.Equal(t, tt.payload["round_id"], scoreData.RoundID)
				assert.Equal(t, tt.payload["winner_team_id"], scoreData.WinnerTeamID)
				assert.Equal(t, tt.payload["loser_team_id"], scoreData.LoserTeamID)
				assert.Equal(t, tt.payload["home_points"], scoreData.HomePoints)
				assert.Equal(t, tt.payload["away_points"], scoreData.AwayPoints)
			}
		})
	}
}

func TestScoreDataJSONTags(t *testing.T) {
	// Test that JSON marshaling/unmarshaling works correctly
	original := ScoreData{
		MatchID:      497,
		RoundID:      494,
		WinnerTeamID: 204,
		LoserTeamID:  268,
		HomePoints:   10,
		AwayPoints:   0, // Test with 0 value
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(original)
	assert.NoError(t, err)

	// Verify JSON contains expected fields
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	assert.NoError(t, err)

	assert.Equal(t, float64(497), jsonMap["match_id"])
	assert.Equal(t, float64(494), jsonMap["round_id"])
	assert.Equal(t, float64(204), jsonMap["winner_team_id"])
	assert.Equal(t, float64(268), jsonMap["loser_team_id"])
	assert.Equal(t, float64(10), jsonMap["home_points"])
	assert.Equal(t, float64(0), jsonMap["away_points"]) // Ensure 0 is preserved

	// Unmarshal back to struct
	var unmarshaled ScoreData
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, original, unmarshaled)
}
