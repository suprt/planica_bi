package repositories

import (
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// DirectRepository handles database operations for Direct data
type DirectRepository struct {
	db *gorm.DB
}

// NewDirectRepository creates a new Direct repository
func NewDirectRepository(db *gorm.DB) *DirectRepository {
	return &DirectRepository{db: db}
}

// GetCampaignMonthly retrieves monthly campaign metrics
func (r *DirectRepository) GetCampaignMonthly(projectID uint, year int, month int) ([]*models.DirectCampaignMonthly, error) {
	var metrics []*models.DirectCampaignMonthly
	err := r.db.Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
		Find(&metrics).Error
	return metrics, err
}

// SaveCampaignMonthly saves campaign monthly metrics
func (r *DirectRepository) SaveCampaignMonthly(metrics *models.DirectCampaignMonthly) error {
	return r.db.Save(metrics).Error
}

// GetTotalsMonthly retrieves monthly totals for a project
func (r *DirectRepository) GetTotalsMonthly(projectID uint, year int, month int) (*models.DirectTotalsMonthly, error) {
	var totals models.DirectTotalsMonthly
	err := r.db.Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
		First(&totals).Error
	if err != nil {
		return nil, err
	}
	return &totals, nil
}

// SaveTotalsMonthly saves monthly totals
func (r *DirectRepository) SaveTotalsMonthly(totals *models.DirectTotalsMonthly) error {
	return r.db.Save(totals).Error
}
