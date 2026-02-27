package repositories

import (
	"context"
	"time"

	"github.com/suprt/planica_bi/backend/internal/cache"
	"github.com/suprt/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// GoalRepository handles database operations for goals
type GoalRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

// NewGoalRepository creates a new goal repository
func NewGoalRepository(db *gorm.DB, cache *cache.Cache) *GoalRepository {
	return &GoalRepository{
		db:    db,
		cache: cache,
	}
}

// Create creates a new goal
func (r *GoalRepository) Create(ctx context.Context, goal *models.Goal) error {
	if err := r.db.WithContext(ctx).Create(goal).Error; err != nil {
		return err
	}

	// Invalidate cache for this counter
	if r.cache != nil {
		cacheKey := cache.BuildKey(cache.KeyPrefixGoals, goal.CounterID)
		_ = r.cache.Delete(cacheKey)
	}

	return nil
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

// GetByCounterID retrieves all goals for a counter (with cache)
func (r *GoalRepository) GetByCounterID(ctx context.Context, counterID uint) ([]*models.Goal, error) {
	// Try to get from cache
	cacheKey := cache.BuildKey(cache.KeyPrefixGoals, counterID)
	var goals []*models.Goal
	if err := r.cache.Get(cacheKey, &goals); err == nil {
		return goals, nil
	}

	// Cache miss - get from database
	err := r.db.WithContext(ctx).Where("counter_id = ?", counterID).Find(&goals).Error
	if err != nil {
		return nil, err
	}

	// Store in cache (TTL: 1 hour)
	if r.cache != nil {
		_ = r.cache.Set(cacheKey, goals, time.Hour)
	}

	return goals, nil
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
	// Get old goal to know which counter's cache to invalidate
	var oldGoal models.Goal
	if err := r.db.WithContext(ctx).First(&oldGoal, goal.ID).Error; err == nil {
		// Invalidate cache for old counter
		if r.cache != nil {
			cacheKey := cache.BuildKey(cache.KeyPrefixGoals, oldGoal.CounterID)
			_ = r.cache.Delete(cacheKey)
		}
	}

	if err := r.db.WithContext(ctx).Save(goal).Error; err != nil {
		return err
	}

	// Invalidate cache for new counter (if counter changed)
	if r.cache != nil && goal.CounterID != oldGoal.CounterID {
		cacheKey := cache.BuildKey(cache.KeyPrefixGoals, goal.CounterID)
		_ = r.cache.Delete(cacheKey)
	}

	return nil
}

// Delete deletes a goal
func (r *GoalRepository) Delete(ctx context.Context, id uint) error {
	// Get goal to know which counter's cache to invalidate
	var goal models.Goal
	if err := r.db.WithContext(ctx).First(&goal, id).Error; err == nil {
		// Invalidate cache for this counter
		if r.cache != nil {
			cacheKey := cache.BuildKey(cache.KeyPrefixGoals, goal.CounterID)
			_ = r.cache.Delete(cacheKey)
		}
	}

	return r.db.WithContext(ctx).Delete(&models.Goal{}, id).Error
}

// DeleteByCounterID deletes all goals for a counter
func (r *GoalRepository) DeleteByCounterID(ctx context.Context, counterID uint) error {
	if err := r.db.WithContext(ctx).Where("counter_id = ?", counterID).Delete(&models.Goal{}).Error; err != nil {
		return err
	}

	// Invalidate cache for this counter
	if r.cache != nil {
		cacheKey := cache.BuildKey(cache.KeyPrefixGoals, counterID)
		_ = r.cache.Delete(cacheKey)
	}

	return nil
}

// GetByCounterIDs retrieves all goals for multiple counters
func (r *GoalRepository) GetByCounterIDs(ctx context.Context, counterIDs []uint) ([]*models.Goal, error) {
	var goals []*models.Goal
	err := r.db.WithContext(ctx).Where("counter_id IN ?", counterIDs).Find(&goals).Error
	return goals, err
}
