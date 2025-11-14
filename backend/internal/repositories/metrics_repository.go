package repositories

import "gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"

// MetricsRepository handles database operations for metrics
type MetricsRepository struct {
	// TODO: add database connection
}

// NewMetricsRepository creates a new metrics repository
func NewMetricsRepository() *MetricsRepository {
	return &MetricsRepository{}
}

// GetMonthlyMetrics retrieves monthly metrics for a project
func (r *MetricsRepository) GetMonthlyMetrics(projectID uint, year int, month int) (*models.MetricsMonthly, error) {
	// TODO: implement
	return nil, nil
}

// SaveMonthlyMetrics saves monthly metrics
func (r *MetricsRepository) SaveMonthlyMetrics(metrics *models.MetricsMonthly) error {
	// TODO: implement
	return nil
}

// GetAgeMetrics retrieves age-based metrics for a project
func (r *MetricsRepository) GetAgeMetrics(projectID uint, year int, month int) ([]*models.MetricsAgeMonthly, error) {
	// TODO: implement
	return nil, nil
}

// SaveAgeMetrics saves age-based metrics
func (r *MetricsRepository) SaveAgeMetrics(metrics *models.MetricsAgeMonthly) error {
	// TODO: implement
	return nil
}
