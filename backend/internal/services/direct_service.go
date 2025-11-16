package services

import (
	"context"
	"errors"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/pkg/utils"
	"gorm.io/gorm"
)

// DirectService handles business logic for Yandex.Direct accounts
type DirectService struct {
	directRepo *repositories.DirectRepository
}

// NewDirectService creates a new Direct service
func NewDirectService(directRepo *repositories.DirectRepository) *DirectService {
	return &DirectService{
		directRepo: directRepo,
	}
}

// CreateAccount creates a new Direct account
func (s *DirectService) CreateAccount(ctx context.Context, account *models.DirectAccount) error {
	// Validate required fields
	if account.ProjectID == 0 {
		return errors.New("project_id is required")
	}
	if account.ClientLogin == "" {
		return errors.New("client_login is required")
	}

	// Check if account with this ClientLogin already exists for this project
	existing, err := s.directRepo.GetAccountByClientLogin(ctx, account.ProjectID, account.ClientLogin)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existing != nil {
		return errors.New("account with this ClientLogin already exists for this project")
	}

	return s.directRepo.CreateAccount(ctx, account)
}

// GetAccountsByProject retrieves all Direct accounts for a project
func (s *DirectService) GetAccountsByProject(ctx context.Context, projectID uint) ([]*models.DirectAccount, error) {
	return s.directRepo.GetAccountsByProjectID(ctx, projectID)
}

// CampaignMetricsRow represents campaign metrics for a month
type CampaignMetricsRow struct {
	Month       string   `json:"month"`
	Impressions int      `json:"impressions"`
	Clicks      int      `json:"clicks"`
	Ctr         float64  `json:"ctr"`
	Cpc         float64  `json:"cpc"`
	Conv        *int     `json:"conv,omitempty"`
	Cpa         *float64 `json:"cpa,omitempty"`
	Cost        float64  `json:"cost"`
}

// CampaignWithMetrics represents a campaign with its metrics
type CampaignWithMetrics struct {
	CampaignID int64               `json:"campaignId"`
	Name       string              `json:"name"`
	Status     string              `json:"status"`
	Rows       []CampaignMetricsRow `json:"rows"`
}

// GetCampaignsWithMetrics retrieves all campaigns for a project with their metrics for the last 3 months
func (s *DirectService) GetCampaignsWithMetrics(ctx context.Context, projectID uint) ([]CampaignWithMetrics, error) {
	// Get all campaigns for the project
	campaigns, err := s.directRepo.GetCampaignsByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if len(campaigns) == 0 {
		return []CampaignWithMetrics{}, nil
	}

	// Prepare periods: M (current), M-1, M-2
	now := time.Now()
	periodData := make([]struct {
		year   int
		month  int
		period string
	}, 3)

	for i := 0; i < 3; i++ {
		periodTime := now.AddDate(0, -i, 0)
		year := periodTime.Year()
		month := int(periodTime.Month())
		periodStr := utils.FormatPeriod(year, month)
		periodData[i] = struct {
			year   int
			month  int
			period string
		}{year: year, month: month, period: periodStr}
	}

	// Build result
	result := make([]CampaignWithMetrics, 0, len(campaigns))

	for _, campaign := range campaigns {
		campaignData := CampaignWithMetrics{
			CampaignID: campaign.CampaignID,
			Name:       campaign.Name,
			Status:     campaign.Status,
			Rows:       []CampaignMetricsRow{},
		}

		// Get metrics for each period
		for _, pd := range periodData {
			metrics, err := s.directRepo.GetCampaignMonthlyByCampaignID(ctx, projectID, campaign.ID, pd.year, pd.month)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}

			if metrics != nil {
				campaignData.Rows = append(campaignData.Rows, CampaignMetricsRow{
					Month:       pd.period,
					Impressions: metrics.Impressions,
					Clicks:      metrics.Clicks,
					Ctr:         metrics.CTRPct,
					Cpc:         metrics.CPC,
					Conv:        metrics.Conversions,
					Cpa:         metrics.CPA,
					Cost:        metrics.Cost,
				})
			}
		}

		result = append(result, campaignData)
	}

	return result, nil
}
