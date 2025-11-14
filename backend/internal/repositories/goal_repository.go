package repositories

import (
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// GoalRepository handles database operations for goals
type GoalRepository struct {
	db *gorm.DB
}

// NewGoalRepository creates a new goal repository
func NewGoalRepository(db *gorm.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

// Create creates a new goal
func (r *GoalRepository) Create(goal *models.Goal) error {
	return r.db.Create(goal).Error
}

// GetByID retrieves a goal by ID
func (r *GoalRepository) GetByID(id uint) (*models.Goal, error) {
	var goal models.Goal
	err := r.db.First(&goal, id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// GetByCounterID retrieves all goals for a counter
func (r *GoalRepository) GetByCounterID(counterID uint) ([]*models.Goal, error) {
	var goals []*models.Goal
	err := r.db.Where("counter_id = ?", counterID).Find(&goals).Error
	return goals, err
}

// GetByGoalID retrieves a goal by Yandex goal ID and counter ID
func (r *GoalRepository) GetByGoalID(counterID uint, goalID int64) (*models.Goal, error) {
	var goal models.Goal
	err := r.db.Where("counter_id = ? AND goal_id = ?", counterID, goalID).First(&goal).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// Update updates a goal
func (r *GoalRepository) Update(goal *models.Goal) error {
	return r.db.Save(goal).Error
}

// Delete deletes a goal
func (r *GoalRepository) Delete(id uint) error {
	return r.db.Delete(&models.Goal{}, id).Error
}

// DeleteByCounterID deletes all goals for a counter
func (r *GoalRepository) DeleteByCounterID(counterID uint) error {
	return r.db.Where("counter_id = ?", counterID).Delete(&models.Goal{}).Error
}

