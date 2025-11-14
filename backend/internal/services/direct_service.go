package services

import (
	"errors"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gorm.io/gorm"
)

// DirectService handles business logic for Yandex.Direct accounts, campaigns and metrics
type DirectService struct {
	directRepo *repositories.DirectRepository
	db         *gorm.DB
}

// NewDirectService creates a new Direct service
func NewDirectService(directRepo *repositories.DirectRepository, db *gorm.DB) *DirectService {
	return &DirectService{
		directRepo: directRepo,
		db:         db,
	}
}

// DirectAccount methods

// CreateAccount creates a new Direct account
func (s *DirectService) CreateAccount(account *models.DirectAccount) error {
	// Check if account with this ClientLogin already exists for this project
	var existing models.DirectAccount
	err := s.db.Where("project_id = ? AND client_login = ?", account.ProjectID, account.ClientLogin).
		First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == nil {
		return errors.New("account with this ClientLogin already exists for this project")
	}

	return s.db.Create(account).Error
}

// GetAccount retrieves a Direct account by ID
func (s *DirectService) GetAccount(id uint) (*models.DirectAccount, error) {
	var account models.DirectAccount
	err := s.db.First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountsByProject retrieves all Direct accounts for a project
func (s *DirectService) GetAccountsByProject(projectID uint) ([]*models.DirectAccount, error) {
	var accounts []*models.DirectAccount
	err := s.db.Where("project_id = ?", projectID).Find(&accounts).Error
	return accounts, err
}

// UpdateAccount updates a Direct account
func (s *DirectService) UpdateAccount(account *models.DirectAccount) error {
	return s.db.Save(account).Error
}

// DeleteAccount deletes a Direct account
func (s *DirectService) DeleteAccount(id uint) error {
	return s.db.Delete(&models.DirectAccount{}, id).Error
}

// DeleteAccountsByProject deletes all Direct accounts for a project
func (s *DirectService) DeleteAccountsByProject(projectID uint) error {
	return s.db.Where("project_id = ?", projectID).Delete(&models.DirectAccount{}).Error
}

// DirectCampaign methods

// CreateCampaign creates a new Direct campaign
func (s *DirectService) CreateCampaign(campaign *models.DirectCampaign) error {
	// Validate that account exists
	var account models.DirectAccount
	err := s.db.First(&account, campaign.DirectAccountID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Direct account not found")
		}
		return err
	}

	// Check if campaign with this CampaignID already exists for this account
	var existing models.DirectCampaign
	err = s.db.Where("direct_account_id = ? AND campaign_id = ?", campaign.DirectAccountID, campaign.CampaignID).
		First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == nil {
		return errors.New("campaign with this CampaignID already exists for this account")
	}

	return s.db.Create(campaign).Error
}

// GetCampaign retrieves a Direct campaign by ID
func (s *DirectService) GetCampaign(id uint) (*models.DirectCampaign, error) {
	var campaign models.DirectCampaign
	err := s.db.First(&campaign, id).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

// GetCampaignsByAccount retrieves all campaigns for a Direct account
func (s *DirectService) GetCampaignsByAccount(accountID uint) ([]*models.DirectCampaign, error) {
	var campaigns []*models.DirectCampaign
	err := s.db.Where("direct_account_id = ?", accountID).Find(&campaigns).Error
	return campaigns, err
}

// GetCampaignsByProject retrieves all campaigns for a project (through accounts)
func (s *DirectService) GetCampaignsByProject(projectID uint) ([]*models.DirectCampaign, error) {
	var campaigns []*models.DirectCampaign
	err := s.db.Table("direct_campaigns").
		Joins("INNER JOIN direct_accounts ON direct_campaigns.direct_account_id = direct_accounts.id").
		Where("direct_accounts.project_id = ?", projectID).
		Find(&campaigns).Error
	return campaigns, err
}

// GetCampaignByCampaignID retrieves a campaign by Yandex CampaignID and account ID
func (s *DirectService) GetCampaignByCampaignID(accountID uint, campaignID int64) (*models.DirectCampaign, error) {
	var campaign models.DirectCampaign
	err := s.db.Where("direct_account_id = ? AND campaign_id = ?", accountID, campaignID).
		First(&campaign).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

// UpdateCampaign updates a Direct campaign
func (s *DirectService) UpdateCampaign(campaign *models.DirectCampaign) error {
	// Validate that account exists if DirectAccountID changed
	if campaign.DirectAccountID != 0 {
		var account models.DirectAccount
		err := s.db.First(&account, campaign.DirectAccountID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("Direct account not found")
			}
			return err
		}
	}

	return s.db.Save(campaign).Error
}

// DeleteCampaign deletes a Direct campaign
func (s *DirectService) DeleteCampaign(id uint) error {
	return s.db.Delete(&models.DirectCampaign{}, id).Error
}

// DeleteCampaignsByAccount deletes all campaigns for a Direct account
func (s *DirectService) DeleteCampaignsByAccount(accountID uint) error {
	return s.db.Where("direct_account_id = ?", accountID).Delete(&models.DirectCampaign{}).Error
}

// Monthly metrics methods

// GetCampaignMonthly retrieves monthly campaign metrics for a project
func (s *DirectService) GetCampaignMonthly(projectID uint, year int, month int) ([]*models.DirectCampaignMonthly, error) {
	return s.directRepo.GetCampaignMonthly(projectID, year, month)
}

// SaveCampaignMonthly saves campaign monthly metrics
func (s *DirectService) SaveCampaignMonthly(metrics *models.DirectCampaignMonthly) error {
	// Check if record exists
	var existing models.DirectCampaignMonthly
	err := s.db.Where("project_id = ? AND direct_campaign_id = ? AND year = ? AND month = ?",
		metrics.ProjectID, metrics.DirectCampaignID, metrics.Year, metrics.Month).
		First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == nil {
		metrics.ID = existing.ID
	}

	return s.directRepo.SaveCampaignMonthly(metrics)
}

// GetTotalsMonthly retrieves monthly totals for a project
func (s *DirectService) GetTotalsMonthly(projectID uint, year int, month int) (*models.DirectTotalsMonthly, error) {
	return s.directRepo.GetTotalsMonthly(projectID, year, month)
}

// SaveTotalsMonthly saves monthly totals
func (s *DirectService) SaveTotalsMonthly(totals *models.DirectTotalsMonthly) error {
	// Check if record exists
	existing, err := s.directRepo.GetTotalsMonthly(totals.ProjectID, totals.Year, totals.Month)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existing != nil {
		totals.ID = existing.ID
	}

	return s.directRepo.SaveTotalsMonthly(totals)
}

// GetCampaignsByAccountID retrieves campaigns by account ID with validation
func (s *DirectService) GetCampaignsByAccountID(accountID uint) ([]*models.DirectCampaign, error) {
	// Validate that account exists
	_, err := s.GetAccount(accountID)
	if err != nil {
		return nil, err
	}

	return s.GetCampaignsByAccount(accountID)
}

