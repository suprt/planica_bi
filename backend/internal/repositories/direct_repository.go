package repositories

import (
	"context"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/cache"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// DirectRepository handles database operations for Direct accounts
type DirectRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

// NewDirectRepository creates a new Direct repository
func NewDirectRepository(db *gorm.DB, cache *cache.Cache) *DirectRepository {
	return &DirectRepository{
		db:    db,
		cache: cache,
	}
}

// CreateAccount creates a new Direct account
func (r *DirectRepository) CreateAccount(ctx context.Context, account *models.DirectAccount) error {
	if err := r.db.WithContext(ctx).Create(account).Error; err != nil {
		return err
	}

	// Invalidate cache for this project
	if r.cache != nil {
		cacheKey := cache.BuildKey(cache.KeyPrefixDirectAccounts, account.ProjectID)
		_ = r.cache.Delete(cacheKey)
	}

	return nil
}

// GetAccountsByProjectID retrieves all Direct accounts for a project (with cache)
func (r *DirectRepository) GetAccountsByProjectID(ctx context.Context, projectID uint) ([]*models.DirectAccount, error) {
	// Try to get from cache
	cacheKey := cache.BuildKey(cache.KeyPrefixDirectAccounts, projectID)
	var accounts []*models.DirectAccount
	if err := r.cache.Get(cacheKey, &accounts); err == nil {
		return accounts, nil
	}

	// Cache miss - get from database
	err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	// Store in cache (TTL: 1 hour)
	if r.cache != nil {
		_ = r.cache.Set(cacheKey, accounts, time.Hour)
	}

	return accounts, err
}

// GetAccountByClientLogin retrieves a Direct account by project ID and client login
func (r *DirectRepository) GetAccountByClientLogin(ctx context.Context, projectID uint, clientLogin string) (*models.DirectAccount, error) {
	var account models.DirectAccount
	err := r.db.WithContext(ctx).Where("project_id = ? AND client_login = ?", projectID, clientLogin).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetCampaignMonthly retrieves monthly campaign metrics
func (r *DirectRepository) GetCampaignMonthly(ctx context.Context, projectID uint, year int, month int) ([]*models.DirectCampaignMonthly, error) {
	var metrics []*models.DirectCampaignMonthly
	err := r.db.WithContext(ctx).Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
		Find(&metrics).Error
	return metrics, err
}

// SaveCampaignMonthly saves campaign monthly metrics
func (r *DirectRepository) SaveCampaignMonthly(metrics *models.DirectCampaignMonthly) error {
	return r.db.Save(metrics).Error
}

// GetTotalsMonthly retrieves monthly totals for a project
func (r *DirectRepository) GetTotalsMonthly(ctx context.Context, projectID uint, year int, month int) (*models.DirectTotalsMonthly, error) {
	var totals models.DirectTotalsMonthly
	err := r.db.WithContext(ctx).Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
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

// GetCampaignsByAccountID retrieves all campaigns for a Direct account (with cache)
func (r *DirectRepository) GetCampaignsByAccountID(ctx context.Context, accountID uint) ([]*models.DirectCampaign, error) {
	// Try to get from cache
	cacheKey := cache.BuildKey(cache.KeyPrefixDirectCampaigns, accountID)
	var campaigns []*models.DirectCampaign
	if err := r.cache.Get(cacheKey, &campaigns); err == nil {
		return campaigns, nil
	}

	// Cache miss - get from database
	err := r.db.WithContext(ctx).Where("direct_account_id = ?", accountID).Find(&campaigns).Error
	if err != nil {
		return nil, err
	}

	// Store in cache (TTL: 1 hour)
	if r.cache != nil {
		_ = r.cache.Set(cacheKey, campaigns, time.Hour)
	}

	return campaigns, err
}

// CreateCampaign creates a new Direct campaign
func (r *DirectRepository) CreateCampaign(ctx context.Context, campaign *models.DirectCampaign) error {
	if err := r.db.WithContext(ctx).Create(campaign).Error; err != nil {
		return err
	}

	// Invalidate cache for this account
	if r.cache != nil {
		cacheKey := cache.BuildKey(cache.KeyPrefixDirectCampaigns, campaign.DirectAccountID)
		_ = r.cache.Delete(cacheKey)
	}

	return nil
}

// GetCampaignByID retrieves a Direct campaign by ID
func (r *DirectRepository) GetCampaignByID(ctx context.Context, id uint) (*models.DirectCampaign, error) {
	var campaign models.DirectCampaign
	err := r.db.WithContext(ctx).First(&campaign, id).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

// GetCampaignMonthlyByCampaignID retrieves campaign monthly metrics by campaign ID
func (r *DirectRepository) GetCampaignMonthlyByCampaignID(ctx context.Context, projectID uint, directCampaignID uint, year int, month int) (*models.DirectCampaignMonthly, error) {
	var metrics models.DirectCampaignMonthly
	err := r.db.WithContext(ctx).Where("project_id = ? AND direct_campaign_id = ? AND year = ? AND month = ?",
		projectID, directCampaignID, year, month).First(&metrics).Error
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}
