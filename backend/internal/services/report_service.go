package services

import (
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/pkg/utils"
	"gorm.io/gorm"
)

// ReportService handles business logic for reports
type ReportService struct {
	metricsRepo *repositories.MetricsRepository
	directRepo  *repositories.DirectRepository
	db          *gorm.DB
}

// NewReportService creates a new report service
func NewReportService(metricsRepo *repositories.MetricsRepository, directRepo *repositories.DirectRepository, db *gorm.DB) *ReportService {
	return &ReportService{
		metricsRepo: metricsRepo,
		directRepo:  directRepo,
		db:          db,
	}
}

// ReportPeriod represents data for a single period (month)
type ReportPeriod struct {
	Period            string                        `json:"period"` // "YYYY-MM"
	Metrica           *models.MetricsMonthly        `json:"metrica,omitempty"`
	AgeBreakdown      []*models.MetricsAgeMonthly   `json:"age_breakdown,omitempty"`
	DirectTotals      *models.DirectTotalsMonthly   `json:"direct_totals,omitempty"`
	DirectCampaigns   []*models.DirectCampaignMonthly `json:"direct_campaigns,omitempty"`
	SEOQueries        []*models.SEOQueriesMonthly   `json:"seo_queries,omitempty"`
}

// Report represents a full report with 3 periods
type Report struct {
	ProjectID uint          `json:"project_id"`
	Periods   []ReportPeriod `json:"periods"`
}

// GetReport generates a report for a project for the last 3 months
func (s *ReportService) GetReport(projectID uint) (*Report, error) {
	now := time.Now()
	
	// Get periods: M (current), M-1, M-2
	periods := make([]ReportPeriod, 3)
	
	for i := 0; i < 3; i++ {
		periodTime := now.AddDate(0, -i, 0)
		year := periodTime.Year()
		month := int(periodTime.Month())
		periodStr := utils.FormatPeriod(year, month)
		
		period := ReportPeriod{
			Period: periodStr,
		}
		
		// Get Metrica summary
		metrics, err := s.metricsRepo.GetMonthlyMetrics(projectID, year, month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if metrics != nil {
			period.Metrica = metrics
		}
		
		// Get age breakdown
		ageMetrics, err := s.metricsRepo.GetAgeMetrics(projectID, year, month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if len(ageMetrics) > 0 {
			period.AgeBreakdown = ageMetrics
		}
		
		// Get Direct totals
		directTotals, err := s.directRepo.GetTotalsMonthly(projectID, year, month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if directTotals != nil {
			period.DirectTotals = directTotals
		}
		
		// Get Direct campaigns
		directCampaigns, err := s.directRepo.GetCampaignMonthly(projectID, year, month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if len(directCampaigns) > 0 {
			period.DirectCampaigns = directCampaigns
		}
		
		// Get SEO queries
		var seoQueries []*models.SEOQueriesMonthly
		err = s.db.Where("project_id = ? AND year = ? AND month = ?", projectID, year, month).
			Find(&seoQueries).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if len(seoQueries) > 0 {
			period.SEOQueries = seoQueries
		}
		
		periods[i] = period
	}
	
	return &Report{
		ProjectID: projectID,
		Periods:   periods,
	}, nil
}

// CalculateDynamics calculates percentage change between two values
func (s *ReportService) CalculateDynamics(current, previous float64) float64 {
	return utils.CalculateDynamics(current, previous)
}

