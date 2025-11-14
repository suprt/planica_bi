package services

import (
	"context"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/pkg/utils"
	"gorm.io/gorm"
)

// ReportService handles business logic for reports
type ReportService struct {
	metricsRepo *repositories.MetricsRepository
	directRepo  *repositories.DirectRepository
	seoRepo     *repositories.SEORepository
}

// NewReportService creates a new report service
func NewReportService(
	metricsRepo *repositories.MetricsRepository,
	directRepo *repositories.DirectRepository,
	seoRepo *repositories.SEORepository,
) *ReportService {
	return &ReportService{
		metricsRepo: metricsRepo,
		directRepo:  directRepo,
		seoRepo:     seoRepo,
	}
}

// MetricaSummaryRow represents a single row in metrica summary
type MetricaSummaryRow struct {
	Month  string  `json:"month"`
	Visits int     `json:"visits"`
	Users  int     `json:"users"`
	Bounce float64 `json:"bounce"`
	AvgSec int     `json:"avgSec"`
	Conv   *int    `json:"conv,omitempty"`
}

// MetricaAgeRow represents a single row in metrica age breakdown
type MetricaAgeRow struct {
	Month  string  `json:"month"`
	Age    string  `json:"age"`
	Visits int     `json:"visits"`
	Users  int     `json:"users"`
	Bounce float64 `json:"bounce"`
	AvgSec int     `json:"avgSec"`
}

// DirectTotalsRow represents a single row in direct totals
type DirectTotalsRow struct {
	Month       string   `json:"month"`
	Impressions int      `json:"impressions"`
	Clicks      int      `json:"clicks"`
	Ctr         float64  `json:"ctr"`
	Cpc         float64  `json:"cpc"`
	Conv        *int     `json:"conv,omitempty"`
	Cpa         *float64 `json:"cpa,omitempty"`
	Cost        float64  `json:"cost"`
}

// DirectCampaignRow represents a single row for a campaign in a month
type DirectCampaignRow struct {
	Month       string   `json:"month"`
	Impressions int      `json:"impressions"`
	Clicks      int      `json:"clicks"`
	Ctr         float64  `json:"ctr"`
	Cpc         float64  `json:"cpc"`
	Conv        *int     `json:"conv,omitempty"`
	Cpa         *float64 `json:"cpa,omitempty"`
	Cost        float64  `json:"cost"`
}

// DirectCampaignData represents campaign data with rows for all months
type DirectCampaignData struct {
	CampaignID int64               `json:"campaignId"`
	Name       string              `json:"name"`
	Rows       []DirectCampaignRow `json:"rows"`
}

// SEOSummaryRow represents a single row in SEO summary
type SEOSummaryRow struct {
	Month    string `json:"month"`
	Visitors int    `json:"visitors"`
	Conv     int    `json:"conv"`
}

// SEOQueryRow represents a single SEO query row
type SEOQueryRow struct {
	Month    string  `json:"month"`
	Query    string  `json:"query"`
	Position int     `json:"position"`
	URL      *string `json:"url,omitempty"`
}

// MetricaData represents metrica section of the report
type MetricaData struct {
	Summary []MetricaSummaryRow `json:"summary"`
	Age     []MetricaAgeRow     `json:"age"`
}

// DirectData represents direct section of the report
type DirectData struct {
	Totals    []DirectTotalsRow    `json:"totals"`
	Campaigns []DirectCampaignData `json:"campaigns"`
}

// SEOData represents SEO section of the report
type SEOData struct {
	Summary []SEOSummaryRow `json:"summary"`
	Queries []SEOQueryRow   `json:"queries"`
}

// Report represents a full report according to TZ format
type Report struct {
	ProjectID uint        `json:"projectId"`
	Periods   []string    `json:"periods"`
	Metrica   MetricaData `json:"metrica"`
	Direct    DirectData  `json:"direct"`
	SEO       SEOData     `json:"seo"`
}

// GetReport generates a report for a project for the last 3 months
func (s *ReportService) GetReport(ctx context.Context, projectID uint) (*Report, error) {
	now := time.Now()

	// Prepare periods array: M (current), M-1, M-2
	periods := make([]string, 3)
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
		periods[i] = periodStr
		periodData[i] = struct {
			year   int
			month  int
			period string
		}{year: year, month: month, period: periodStr}
	}

	report := &Report{
		ProjectID: projectID,
		Periods:   periods,
		Metrica: MetricaData{
			Summary: []MetricaSummaryRow{},
			Age:     []MetricaAgeRow{},
		},
		Direct: DirectData{
			Totals:    []DirectTotalsRow{},
			Campaigns: []DirectCampaignData{},
		},
		SEO: SEOData{
			Summary: []SEOSummaryRow{},
			Queries: []SEOQueryRow{},
		},
	}

	// Map to group campaigns by CampaignID (Yandex ID)
	campaignMap := make(map[int64]*DirectCampaignData)

	// Process each period
	for _, pd := range periodData {
		// Get Metrica summary
		metrics, err := s.metricsRepo.GetMonthlyMetrics(ctx, projectID, pd.year, pd.month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if metrics != nil {
			report.Metrica.Summary = append(report.Metrica.Summary, MetricaSummaryRow{
				Month:  pd.period,
				Visits: metrics.Visits,
				Users:  metrics.Users,
				Bounce: metrics.BounceRate,
				AvgSec: metrics.AvgSessionDurationSec,
				Conv:   metrics.Conversions,
			})
		}

		// Get age breakdown
		ageMetrics, err := s.metricsRepo.GetAgeMetrics(ctx, projectID, pd.year, pd.month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		for _, age := range ageMetrics {
			report.Metrica.Age = append(report.Metrica.Age, MetricaAgeRow{
				Month:  pd.period,
				Age:    string(age.AgeGroup),
				Visits: age.Visits,
				Users:  age.Users,
				Bounce: age.BounceRate,
				AvgSec: age.AvgSessionDurationSec,
			})
		}

		// Get Direct totals
		directTotals, err := s.directRepo.GetTotalsMonthly(ctx, projectID, pd.year, pd.month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if directTotals != nil {
			report.Direct.Totals = append(report.Direct.Totals, DirectTotalsRow{
				Month:       pd.period,
				Impressions: directTotals.Impressions,
				Clicks:      directTotals.Clicks,
				Ctr:         directTotals.CTRPct,
				Cpc:         directTotals.CPC,
				Conv:        directTotals.Conversions,
				Cpa:         directTotals.CPA,
				Cost:        directTotals.Cost,
			})
		}

		// Get Direct campaigns and group by CampaignID
		directCampaigns, err := s.directRepo.GetCampaignMonthly(ctx, projectID, pd.year, pd.month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		for _, campaignMonthly := range directCampaigns {
			// Get campaign info
			campaign, err := s.directRepo.GetCampaignByID(ctx, campaignMonthly.DirectCampaignID)
			if err != nil {
				continue // Skip if campaign not found
			}

			// Get or create campaign data
			if _, exists := campaignMap[campaign.CampaignID]; !exists {
				campaignMap[campaign.CampaignID] = &DirectCampaignData{
					CampaignID: campaign.CampaignID,
					Name:       campaign.Name,
					Rows:       []DirectCampaignRow{},
				}
			}

			// Add row for this month
			campaignMap[campaign.CampaignID].Rows = append(campaignMap[campaign.CampaignID].Rows, DirectCampaignRow{
				Month:       pd.period,
				Impressions: campaignMonthly.Impressions,
				Clicks:      campaignMonthly.Clicks,
				Ctr:         campaignMonthly.CTRPct,
				Cpc:         campaignMonthly.CPC,
				Conv:        campaignMonthly.Conversions,
				Cpa:         campaignMonthly.CPA,
				Cost:        campaignMonthly.Cost,
			})
		}

		// Get SEO queries
		seoQueries, err := s.seoRepo.GetSEOQueries(ctx, projectID, pd.year, pd.month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		for _, query := range seoQueries {
			report.SEO.Queries = append(report.SEO.Queries, SEOQueryRow{
				Month:    pd.period,
				Query:    query.Query,
				Position: query.Position,
				URL:      query.URL,
			})
		}
	}

	// Convert campaign map to slice
	for _, campaignData := range campaignMap {
		report.Direct.Campaigns = append(report.Direct.Campaigns, *campaignData)
	}

	// TODO: SEO summary - need to get organic visitors and conversions from Metrica
	// For now, leaving empty as it requires Metrica API integration with organic segment

	return report, nil
}

// CalculateDynamics calculates percentage change between two values
func (s *ReportService) CalculateDynamics(current, previous float64) float64 {
	return utils.CalculateDynamics(current, previous)
}
