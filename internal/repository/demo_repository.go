package repository

import (
	"fmt"

	"github.com/UDL-TF/UnitedAPI/internal/model"
	"gorm.io/gorm"
)

// DemoRepository handles demo database operations
type DemoRepository struct {
	db *gorm.DB
}

// NewDemoRepository creates a new demo repository
func NewDemoRepository(db *gorm.DB) *DemoRepository {
	return &DemoRepository{db: db}
}

// Create saves a new demo record to the database
func (r *DemoRepository) Create(demo *model.Demo) error {
	if err := r.db.Create(demo).Error; err != nil {
		return fmt.Errorf("failed to create demo record: %w", err)
	}
	return nil
}

// GetByID retrieves a demo by its ID
func (r *DemoRepository) GetByID(id int) (*model.Demo, error) {
	var demo model.Demo
	if err := r.db.First(&demo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("demo not found with ID %d", id)
		}
		return nil, fmt.Errorf("failed to get demo by ID: %w", err)
	}
	return &demo, nil
}

// GetByMatchAndRound retrieves a demo by match ID and round ID
func (r *DemoRepository) GetByMatchAndRound(matchID, roundID int) (*model.Demo, error) {
	var demo model.Demo
	if err := r.db.Where("match_id = ? AND round_id = ?", matchID, roundID).First(&demo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("demo not found for match %d, round %d", matchID, roundID)
		}
		return nil, fmt.Errorf("failed to get demo by match and round: %w", err)
	}
	return &demo, nil
}

// GetByTournament retrieves demos by tournament ID
func (r *DemoRepository) GetByTournament(tournamentID int) ([]*model.Demo, error) {
	var demos []*model.Demo
	if err := r.db.Where("tournament_id = ?", tournamentID).Find(&demos).Error; err != nil {
		return nil, fmt.Errorf("failed to get demos by tournament: %w", err)
	}
	return demos, nil
}

// GetByTournamentMatchRound retrieves a demo by tournament, match, and round IDs
func (r *DemoRepository) GetByTournamentMatchRound(tournamentID, matchID, roundID int) (*model.Demo, error) {
	var demo model.Demo
	if err := r.db.Where("tournament_id = ? AND match_id = ? AND round_id = ?",
		tournamentID, matchID, roundID).First(&demo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("demo not found for tournament %d, match %d, round %d",
				tournamentID, matchID, roundID)
		}
		return nil, fmt.Errorf("failed to get demo by tournament, match and round: %w", err)
	}
	return &demo, nil
}

// List retrieves all demos with pagination
func (r *DemoRepository) List(offset, limit int) ([]*model.Demo, error) {
	var demos []*model.Demo
	if err := r.db.Offset(offset).Limit(limit).Order("uploaded_at DESC").Find(&demos).Error; err != nil {
		return nil, fmt.Errorf("failed to list demos: %w", err)
	}
	return demos, nil
}

// Delete removes a demo record from the database
func (r *DemoRepository) Delete(id int) error {
	if err := r.db.Delete(&model.Demo{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete demo: %w", err)
	}
	return nil
}

// ExistsByObjectName checks if a demo with the given object name exists
func (r *DemoRepository) ExistsByObjectName(objectName string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Demo{}).Where("object_name = ?", objectName).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check demo existence: %w", err)
	}
	return count > 0, nil
}
