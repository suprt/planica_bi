package services

import (
	"context"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/pkg/utils"
	"gorm.io/gorm"
)

// MetricsService handles business logic for Yandex.Metrica metrics
type MetricsService struct {
	metricsRepo *repositories.MetricsRepository
}

// NewMetricsService creates a new Metrics service
func NewMetricsService(metricsRepo *repositories.MetricsRepository) *MetricsService {
	return &MetricsService{
		metricsRepo: metricsRepo,
	}
}

// MetricsRow represents metrics for a month
type MetricsRow struct {
	Month       string  `json:"month"`
	Visits      int     `json:"visits"`
	Users       int     `json:"users"`
	BounceRate  float64 `json:"bounce_rate"`
	AvgSec      int     `json:"avg_sec"`
	Conversions *int    `json:"conversions,omitempty"`
}

// MetricsWithData represents metrics data for a project
type MetricsWithData struct {
	ProjectID uint        `json:"projectId"`
	Rows      []MetricsRow `json:"rows"`
}

// GetMetricsWithData retrieves metrics for a project for the last 3 months
func (s *MetricsService) GetMetricsWithData(ctx context.Context, projectID uint) (*MetricsWithData, error) {
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

	result := &MetricsWithData{
		ProjectID: projectID,
		Rows:      []MetricsRow{},
	}

	// Get metrics for each period
	for _, pd := range periodData {
		metrics, err := s.metricsRepo.GetMonthlyMetrics(ctx, projectID, pd.year, pd.month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if metrics != nil {
			result.Rows = append(result.Rows, MetricsRow{
				Month:       pd.period,
				Visits:      metrics.Visits,
				Users:       metrics.Users,
				BounceRate:  metrics.BounceRate,
				AvgSec:      metrics.AvgSessionDurationSec,
				Conversions: metrics.Conversions,
			})
		}
	}

	return result, nil
}

