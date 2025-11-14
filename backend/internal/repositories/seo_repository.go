package repositories

import (
	"context"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// SEORepository handles database operations for SEO queries
type SEORepository struct {
	db *gorm.DB
}

// NewSEORepository creates a new SEO repository
func NewSEORepository(db *gorm.DB) *SEORepository {
	return &SEORepository{db: db}
}

// GetSEOQueries retrieves SEO queries for a project, year and month
func (r *SEORepository) GetSEOQueries(ctx context.Context, projectID uint, year int, month int) ([]*models.SEOQueriesMonthly, error) {
	var queries []*models.SEOQueriesMonthly
	err := r.db.WithContext(ctx).Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
		Find(&queries).Error
	return queries, err
}
