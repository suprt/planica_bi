package repositories

import (
	"context"

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
func (r *GoalRepository) Create(ctx context.Context, goal *models.Goal) error {
	return r.db.WithContext(ctx).Create(goal).Error
}

// GetByID retrieves a goal by ID
func (r *GoalRepository) GetByID(ctx context.Context, id uint) (*models.Goal, error) {
	var goal models.Goal
	err := r.db.WithContext(ctx).First(&goal, id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// GetByCounterID retrieves all goals for a counter
func (r *GoalRepository) GetByCounterID(ctx context.Context, counterID uint) ([]*models.Goal, error) {
	var goals []*models.Goal
	err := r.db.WithContext(ctx).Where("counter_id = ?", counterID).Find(&goals).Error
	return goals, err
}

// GetByGoalID retrieves a goal by Yandex goal ID and counter ID
func (r *GoalRepository) GetByGoalID(ctx context.Context, counterID uint, goalID int64) (*models.Goal, error) {
	var goal models.Goal
	err := r.db.WithContext(ctx).Where("counter_id = ? AND goal_id = ?", counterID, goalID).First(&goal).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// Update updates a goal
func (r *GoalRepository) Update(ctx context.Context, goal *models.Goal) error {
	return r.db.WithContext(ctx).Save(goal).Error
}

// Delete deletes a goal
func (r *GoalRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Goal{}, id).Error
}

// DeleteByCounterID deletes all goals for a counter
func (r *GoalRepository) DeleteByCounterID(ctx context.Context, counterID uint) error {
	return r.db.WithContext(ctx).Where("counter_id = ?", counterID).Delete(&models.Goal{}).Error
}

// GetByCounterIDs retrieves all goals for multiple counters
func (r *GoalRepository) GetByCounterIDs(ctx context.Context, counterIDs []uint) ([]*models.Goal, error) {
	var goals []*models.Goal
	err := r.db.WithContext(ctx).Where("counter_id IN ?", counterIDs).Find(&goals).Error
	return goals, err
}
