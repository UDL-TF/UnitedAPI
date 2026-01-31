package handler

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/UDL-TF/UnitedAPI/internal/response"
	"github.com/gin-gonic/gin"
)

// GetDemo handles GET /api/v1/demos
// Query params: match_id, round_id
// Returns the demo file for the given match and round
func GetDemo(c *gin.Context) {
	matchID := c.Query("match_id")
	roundID := c.Query("round_id")

	if matchID == "" || roundID == "" {
		response.BadRequest(c, "Missing match_id or round_id")
		return
	}

	// Force matchID and roundID to be numbers for security
	matchIDNum, err := strconv.Atoi(matchID)
	if err != nil {
		response.BadRequest(c, "Match ID is not a number")
		return
	}

	roundIDNum, err := strconv.Atoi(roundID)
	if err != nil {
		response.BadRequest(c, "Round ID is not a number")
		return
	}

	// Get demo path from environment variable or use default
	uploadPath := os.Getenv("DEMO_PATH")
	if uploadPath == "" {
		uploadPath = "./demo"
	}

	// Find the file with the given match ID and round ID using regex
	filePattern := fmt.Sprintf(`match-%d-round-%d`, matchIDNum, roundIDNum)
	files, err := os.ReadDir(uploadPath)
	if err != nil {
		response.InternalServerError(c, "Error reading demo directory")
		return
	}

	var demoFile string
	var maxSize int64

	for _, file := range files {
		if !file.IsDir() && matchRegex(file.Name(), filePattern) {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}
			// Get the largest matching file
			if fileInfo.Size() > maxSize {
				maxSize = fileInfo.Size()
				demoFile = file.Name()
			}
		}
	}

	if demoFile == "" {
		response.NotFound(c, "Demo file not found")
		return
	}

	fmt.Printf("Downloading demo file: %s\n", demoFile)

	// Set the original file name in the response header
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", demoFile))
	c.File(fmt.Sprintf("%s/%s", uploadPath, demoFile))
}

// matchRegex checks if a string matches a regex pattern
func matchRegex(str, pattern string) bool {
	pattern = regexp.QuoteMeta(pattern)
	matched, _ := regexp.MatchString(pattern, str)
	return matched
}
