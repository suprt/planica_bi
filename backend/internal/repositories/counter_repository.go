package repositories

import (
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
func (r *CounterRepository) Create(counter *models.YandexCounter) error {
	return r.db.Create(counter).Error
}

// GetByID retrieves a counter by ID
func (r *CounterRepository) GetByID(id uint) (*models.YandexCounter, error) {
	var counter models.YandexCounter
	err := r.db.First(&counter, id).Error
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

// GetByProjectID retrieves all counters for a project
func (r *CounterRepository) GetByProjectID(projectID uint) ([]*models.YandexCounter, error) {
	var counters []*models.YandexCounter
	err := r.db.Where("project_id = ?", projectID).Find(&counters).Error
	return counters, err
}

// GetByCounterID retrieves a counter by Yandex counter ID
func (r *CounterRepository) GetByCounterID(counterID int64) (*models.YandexCounter, error) {
	var counter models.YandexCounter
	err := r.db.Where("counter_id = ?", counterID).First(&counter).Error
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

// GetPrimaryByProjectID retrieves the primary counter for a project
func (r *CounterRepository) GetPrimaryByProjectID(projectID uint) (*models.YandexCounter, error) {
	var counter models.YandexCounter
	err := r.db.Where("project_id = ? AND is_primary = ?", projectID, true).First(&counter).Error
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

// Update updates a counter
func (r *CounterRepository) Update(counter *models.YandexCounter) error {
	return r.db.Save(counter).Error
}

// Delete deletes a counter
func (r *CounterRepository) Delete(id uint) error {
	return r.db.Delete(&models.YandexCounter{}, id).Error
}

// DeleteByProjectID deletes all counters for a project
func (r *CounterRepository) DeleteByProjectID(projectID uint) error {
	return r.db.Where("project_id = ?", projectID).Delete(&models.YandexCounter{}).Error
}

