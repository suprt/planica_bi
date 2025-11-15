package repositories

import (
	"context"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/cache"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// CounterRepository handles database operations for Yandex counters
type CounterRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

// NewCounterRepository creates a new counter repository
func NewCounterRepository(db *gorm.DB, cache *cache.Cache) *CounterRepository {
	return &CounterRepository{
		db:    db,
		cache: cache,
	}
}

// Create creates a new counter
func (r *CounterRepository) Create(ctx context.Context, counter *models.YandexCounter) error {
	if err := r.db.WithContext(ctx).Create(counter).Error; err != nil {
		return err
	}

	// Invalidate cache for this project
	if r.cache != nil {
		cacheKey := cache.BuildKey(cache.KeyPrefixCounters, counter.ProjectID)
		_ = r.cache.Delete(cacheKey)
	}

	return nil
}

// GetByID retrieves a counter by ID
func (r *CounterRepository) GetByID(ctx context.Context, id uint) (*models.YandexCounter, error) {
	var counter models.YandexCounter
	err := r.db.WithContext(ctx).First(&counter, id).Error
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

// GetByProjectID retrieves all counters for a project (with cache)
// Counters are cached as they only change on manual admin actions, not during sync
func (r *CounterRepository) GetByProjectID(ctx context.Context, projectID uint) ([]*models.YandexCounter, error) {
	// Try to get from cache
	cacheKey := cache.BuildKey(cache.KeyPrefixCounters, projectID)
	var counters []*models.YandexCounter
	if err := r.cache.Get(cacheKey, &counters); err == nil {
		// Cache hit
		return counters, nil
	}

	// Cache miss - get from database
	err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&counters).Error
	if err != nil {
		return nil, err
	}

	// Store in cache (TTL: 1 hour)
	if r.cache != nil {
		_ = r.cache.Set(cacheKey, counters, time.Hour)
	}

	return counters, nil
}

// GetByCounterID retrieves a counter by Yandex counter ID
func (r *CounterRepository) GetByCounterID(ctx context.Context, counterID int64) (*models.YandexCounter, error) {
	var counter models.YandexCounter
	err := r.db.WithContext(ctx).Where("counter_id = ?", counterID).First(&counter).Error
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

// Update updates a counter
func (r *CounterRepository) Update(ctx context.Context, counter *models.YandexCounter) error {
	// Get old counter to know which project's cache to invalidate
	var oldCounter models.YandexCounter
	if err := r.db.WithContext(ctx).First(&oldCounter, counter.ID).Error; err == nil {
		// Invalidate cache for old project
		if r.cache != nil {
			cacheKey := cache.BuildKey(cache.KeyPrefixCounters, oldCounter.ProjectID)
			_ = r.cache.Delete(cacheKey)
		}
	}

	if err := r.db.WithContext(ctx).Save(counter).Error; err != nil {
		return err
	}

	// Invalidate cache for new project (if project changed)
	if r.cache != nil && counter.ProjectID != oldCounter.ProjectID {
		cacheKey := cache.BuildKey(cache.KeyPrefixCounters, counter.ProjectID)
		_ = r.cache.Delete(cacheKey)
	}

	return nil
}
