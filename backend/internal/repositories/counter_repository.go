package repositories

import (
	"context"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// CounterRepository handles database operations for Yandex counters
type CounterRepository struct {
	db *gorm.DB
}

// NewCounterRepository creates a new counter repository
func NewCounterRepository(db *gorm.DB) *CounterRepository {
	return &CounterRepository{db: db}
}

// Create creates a new counter
func (r *CounterRepository) Create(ctx context.Context, counter *models.YandexCounter) error {
	return r.db.WithContext(ctx).Create(counter).Error
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

// GetByProjectID retrieves all counters for a project
func (r *CounterRepository) GetByProjectID(ctx context.Context, projectID uint) ([]*models.YandexCounter, error) {
	var counters []*models.YandexCounter
	err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&counters).Error
	return counters, err
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
	return r.db.WithContext(ctx).Save(counter).Error
}
