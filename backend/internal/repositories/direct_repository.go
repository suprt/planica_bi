package repositories

import "gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"

// DirectRepository handles database operations for Direct data
type DirectRepository struct {
	// TODO: add database connection
}

// NewDirectRepository creates a new Direct repository
func NewDirectRepository() *DirectRepository {
	return &DirectRepository{}
}

// GetCampaignMonthly retrieves monthly campaign metrics
func (r *DirectRepository) GetCampaignMonthly(projectID uint, year int, month int) ([]*models.DirectCampaignMonthly, error) {
	// TODO: implement
	return nil, nil
}

// SaveCampaignMonthly saves campaign monthly metrics
func (r *DirectRepository) SaveCampaignMonthly(metrics *models.DirectCampaignMonthly) error {
	// TODO: implement
	return nil
}

// GetTotalsMonthly retrieves monthly totals for a project
func (r *DirectRepository) GetTotalsMonthly(projectID uint, year int, month int) (*models.DirectTotalsMonthly, error) {
	// TODO: implement
	return nil, nil
}

// SaveTotalsMonthly saves monthly totals
func (r *DirectRepository) SaveTotalsMonthly(totals *models.DirectTotalsMonthly) error {
	// TODO: implement
	return nil
}
