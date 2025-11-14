package repositories

import (
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// MetricsRepository handles database operations for metrics
type MetricsRepository struct {
	db *gorm.DB
}

// NewMetricsRepository creates a new metrics repository
func NewMetricsRepository(db *gorm.DB) *MetricsRepository {
	return &MetricsRepository{db: db}
}

// GetMonthlyMetrics retrieves monthly metrics for a project
func (r *MetricsRepository) GetMonthlyMetrics(projectID uint, year int, month int) (*models.MetricsMonthly, error) {
	var metrics models.MetricsMonthly
	err := r.db.Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
		First(&metrics).Error
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

// SaveMonthlyMetrics saves monthly metrics
func (r *MetricsRepository) SaveMonthlyMetrics(metrics *models.MetricsMonthly) error {
	return r.db.Save(metrics).Error
}

// GetAgeMetrics retrieves age-based metrics for a project
func (r *MetricsRepository) GetAgeMetrics(projectID uint, year int, month int) ([]*models.MetricsAgeMonthly, error) {
	var metrics []*models.MetricsAgeMonthly
	err := r.db.Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
		Find(&metrics).Error
	return metrics, err
}

// SaveAgeMetrics saves age-based metrics
func (r *MetricsRepository) SaveAgeMetrics(metrics *models.MetricsAgeMonthly) error {
	return r.db.Save(metrics).Error
}
