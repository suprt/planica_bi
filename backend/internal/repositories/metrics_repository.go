package repositories

import (
	"context"

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
func (r *MetricsRepository) GetMonthlyMetrics(ctx context.Context, projectID uint, year int, month int) (*models.MetricsMonthly, error) {
	var metrics models.MetricsMonthly
	err := r.db.WithContext(ctx).Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
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
func (r *MetricsRepository) GetAgeMetrics(ctx context.Context, projectID uint, year int, month int) ([]*models.MetricsAgeMonthly, error) {
	var metrics []*models.MetricsAgeMonthly
	err := r.db.WithContext(ctx).Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
		Find(&metrics).Error
	return metrics, err
}

// SaveAgeMetrics saves age-based metrics
func (r *MetricsRepository) SaveAgeMetrics(metrics *models.MetricsAgeMonthly) error {
	return r.db.Save(metrics).Error
}

// GetAgeMetricsByGroup retrieves age metrics by project, year, month and age group
func (r *MetricsRepository) GetAgeMetricsByGroup(ctx context.Context, projectID uint, year int, month int, ageGroup string) (*models.MetricsAgeMonthly, error) {
	var metrics models.MetricsAgeMonthly
	err := r.db.WithContext(ctx).Where("project_id = ? AND year = ? AND month = ? AND age_group = ?",
		projectID, year, month, ageGroup).First(&metrics).Error
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}
