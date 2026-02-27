package services

import (
	"context"
	"fmt"
	"time"

	"github.com/suprt/planica_bi/backend/internal/models"
)

// MarketingService handles business logic for marketing data
type MarketingService struct {
	directRepo DirectRepositoryInterface
}

// NewMarketingService creates a new marketing service
func NewMarketingService(directRepo DirectRepositoryInterface) *MarketingService {
	return &MarketingService{
		directRepo: directRepo,
	}
}

// MarketingData represents marketing data with clicks and conversions
type MarketingData struct {
	Clicks      MarketingSection `json:"clicks"`
	Conversions MarketingSection `json:"conversions"`
}

// MarketingSection represents a section (clicks or conversions) with summary and metrics
type MarketingSection struct {
	Summary []SummaryItem `json:"summary"`
	Metrics []MetricItem  `json:"metrics"`
}

// SummaryItem represents a summary item with label, value, change and isPositive flag
type SummaryItem struct {
	Label      string  `json:"label"`
	Value      string  `json:"value"`
	Change     float64 `json:"change"`
	IsPositive bool    `json:"isPositive"`
}

// MetricItem represents a metric item with indicator and monthly values
type MetricItem struct {
	ID         int         `json:"id"`
	Indicator  string      `json:"indicator"`
	October    interface{} `json:"october"`
	September  interface{} `json:"september"`
	August     interface{} `json:"august"`
	Efficiency float64     `json:"efficiency"`
}

// GetMarketingData retrieves marketing data for a project for the last 3 months
func (s *MarketingService) GetMarketingData(ctx context.Context, projectID uint) (*MarketingData, error) {
	now := time.Now()

	// Get last 3 months (current, -1, -2)
	currentYear, currentMonth, _ := now.Date()

	months := []struct {
		year  int
		month int
		label string
	}{
		{currentYear, int(currentMonth), "october"},
		{currentYear, int(currentMonth) - 1, "september"},
		{currentYear, int(currentMonth) - 2, "august"},
	}

	// Handle year overflow
	for i := range months {
		if months[i].month <= 0 {
			months[i].month += 12
			months[i].year--
		}
	}

	// Get data for each month
	var october, september, august *models.DirectTotalsMonthly
	var err error

	october, err = s.directRepo.GetTotalsMonthly(ctx, projectID, months[0].year, months[0].month)
	if err != nil {
		return nil, fmt.Errorf("failed to get october data: %w", err)
	}

	september, err = s.directRepo.GetTotalsMonthly(ctx, projectID, months[1].year, months[1].month)
	if err != nil {
		return nil, fmt.Errorf("failed to get september data: %w", err)
	}

	august, err = s.directRepo.GetTotalsMonthly(ctx, projectID, months[2].year, months[2].month)
	if err != nil {
		return nil, fmt.Errorf("failed to get august data: %w", err)
	}

	// Build marketing data
	data := &MarketingData{
		Clicks:      s.buildClicksSection(october, september, august),
		Conversions: s.buildConversionsSection(october, september, august),
	}

	return data, nil
}

// buildClicksSection builds clicks section with summary and metrics
func (s *MarketingService) buildClicksSection(oct, sept, aug *models.DirectTotalsMonthly) MarketingSection {
	var summary []SummaryItem
	var metrics []MetricItem

	// Calculate values with defaults
	octClicks := 0
	septClicks := 0
	augClicks := 0
	octCTR := 0.0
	septCTR := 0.0
	augCTR := 0.0

	if oct != nil {
		octClicks = oct.Clicks
		octCTR = oct.CTRPct
	}
	if sept != nil {
		septClicks = sept.Clicks
		septCTR = sept.CTRPct
	}
	if aug != nil {
		augClicks = aug.Clicks
		augCTR = aug.CTRPct
	}

	// Calculate dynamics
	clicksChange := s.calculateChange(float64(octClicks), float64(septClicks))
	ctrChange := s.calculateChange(octCTR, septCTR)

	// Build summary for clicks
	summary = append(summary, SummaryItem{
		Label:      "Клики",
		Value:      s.formatChangeText(clicksChange),
		Change:     clicksChange,
		IsPositive: clicksChange > 0,
	})

	// For now, we don't have RSA/MC split in the model, so we'll use totals
	// TODO: Add RSA/MC split if needed
	summary = append(summary, SummaryItem{
		Label:      "Клики в RSA",
		Value:      s.formatChangeText(clicksChange),
		Change:     clicksChange,
		IsPositive: clicksChange > 0,
	})

	summary = append(summary, SummaryItem{
		Label:      "Клики в MC",
		Value:      s.formatChangeText(clicksChange),
		Change:     clicksChange,
		IsPositive: clicksChange > 0,
	})

	summary = append(summary, SummaryItem{
		Label:      "CTR в RSA",
		Value:      s.formatChangeText(ctrChange),
		Change:     ctrChange,
		IsPositive: ctrChange > 0,
	})

	summary = append(summary, SummaryItem{
		Label:      "CTR в MC",
		Value:      s.formatChangeText(ctrChange),
		Change:     ctrChange,
		IsPositive: ctrChange > 0,
	})

	// Build metrics for clicks
	metrics = append(metrics, MetricItem{
		ID:         1,
		Indicator:  "Клики, кол-во",
		October:    octClicks,
		September:  septClicks,
		August:     augClicks,
		Efficiency: s.calculateChange(float64(octClicks), float64(septClicks)),
	})

	metrics = append(metrics, MetricItem{
		ID:         2,
		Indicator:  "Клики MC, кол-во",
		October:    octClicks,
		September:  septClicks,
		August:     augClicks,
		Efficiency: s.calculateChange(float64(octClicks), float64(septClicks)),
	})

	metrics = append(metrics, MetricItem{
		ID:         3,
		Indicator:  "Клики RSA, кол-во",
		October:    octClicks,
		September:  septClicks,
		August:     augClicks,
		Efficiency: s.calculateChange(float64(octClicks), float64(septClicks)),
	})

	metrics = append(metrics, MetricItem{
		ID:         4,
		Indicator:  "CTR, %",
		October:    fmt.Sprintf("%.2f%%", octCTR),
		September:  fmt.Sprintf("%.2f%%", septCTR),
		August:     fmt.Sprintf("%.2f%%", augCTR),
		Efficiency: ctrChange,
	})

	metrics = append(metrics, MetricItem{
		ID:         5,
		Indicator:  "CTR MC, %",
		October:    fmt.Sprintf("%.2f%%", octCTR),
		September:  fmt.Sprintf("%.2f%%", septCTR),
		August:     fmt.Sprintf("%.2f%%", augCTR),
		Efficiency: ctrChange,
	})

	metrics = append(metrics, MetricItem{
		ID:         6,
		Indicator:  "CTR RSA, %",
		October:    fmt.Sprintf("%.2f%%", octCTR),
		September:  fmt.Sprintf("%.2f%%", septCTR),
		August:     fmt.Sprintf("%.2f%%", augCTR),
		Efficiency: ctrChange,
	})

	return MarketingSection{
		Summary: summary,
		Metrics: metrics,
	}
}

// buildConversionsSection builds conversions section with summary and metrics
func (s *MarketingService) buildConversionsSection(oct, sept, aug *models.DirectTotalsMonthly) MarketingSection {
	var summary []SummaryItem
	var metrics []MetricItem

	// Calculate values with defaults
	octConversions := 0
	septConversions := 0
	augConversions := 0
	octCPA := 0.0
	septCPA := 0.0
	augCPA := 0.0

	if oct != nil && oct.Conversions != nil {
		octConversions = *oct.Conversions
		if oct.CPA != nil {
			octCPA = *oct.CPA
		}
	}
	if sept != nil && sept.Conversions != nil {
		septConversions = *sept.Conversions
		if sept.CPA != nil {
			septCPA = *sept.CPA
		}
	}
	if aug != nil && aug.Conversions != nil {
		augConversions = *aug.Conversions
		if aug.CPA != nil {
			augCPA = *aug.CPA
		}
	}

	// Calculate dynamics
	conversionsChange := s.calculateChange(float64(octConversions), float64(septConversions))
	cpaChange := s.calculateChange(octCPA, septCPA)

	// Build summary for conversions
	summary = append(summary, SummaryItem{
		Label:      "Конверсии",
		Value:      s.formatChangeText(conversionsChange),
		Change:     conversionsChange,
		IsPositive: conversionsChange > 0,
	})

	summary = append(summary, SummaryItem{
		Label:      "Конверсии в RSA",
		Value:      s.formatChangeText(conversionsChange),
		Change:     conversionsChange,
		IsPositive: conversionsChange > 0,
	})

	summary = append(summary, SummaryItem{
		Label:      "Конверсии в MC",
		Value:      s.formatChangeText(conversionsChange),
		Change:     conversionsChange,
		IsPositive: conversionsChange > 0,
	})

	summary = append(summary, SummaryItem{
		Label:      "CPA",
		Value:      s.formatChangeText(cpaChange),
		Change:     cpaChange,
		IsPositive: false, // CPA increase is bad
	})

	summary = append(summary, SummaryItem{
		Label:      "CPA в MC",
		Value:      s.formatChangeText(cpaChange),
		Change:     cpaChange,
		IsPositive: false, // CPA increase is bad
	})

	// Build metrics for conversions
	metrics = append(metrics, MetricItem{
		ID:         1,
		Indicator:  "Конверсии, кол-во",
		October:    octConversions,
		September:  septConversions,
		August:     augConversions,
		Efficiency: conversionsChange,
	})

	metrics = append(metrics, MetricItem{
		ID:         2,
		Indicator:  "Конверсии MC, кол-во",
		October:    octConversions,
		September:  septConversions,
		August:     augConversions,
		Efficiency: conversionsChange,
	})

	metrics = append(metrics, MetricItem{
		ID:         3,
		Indicator:  "Конверсии RSA, кол-во",
		October:    octConversions,
		September:  septConversions,
		August:     augConversions,
		Efficiency: conversionsChange,
	})

	metrics = append(metrics, MetricItem{
		ID:         4,
		Indicator:  "CPA (средняя цена конверсии), руб",
		October:    octCPA,
		September:  septCPA,
		August:     augCPA,
		Efficiency: cpaChange,
	})

	metrics = append(metrics, MetricItem{
		ID:         5,
		Indicator:  "CPA MC, руб",
		October:    octCPA,
		September:  septCPA,
		August:     augCPA,
		Efficiency: cpaChange,
	})

	metrics = append(metrics, MetricItem{
		ID:         6,
		Indicator:  "CPA RSA, руб",
		October:    octCPA,
		September:  septCPA,
		August:     augCPA,
		Efficiency: cpaChange,
	})

	return MarketingSection{
		Summary: summary,
		Metrics: metrics,
	}
}

// calculateChange calculates percentage change between two values
func (s *MarketingService) calculateChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100 // 100% increase from 0
	}
	return ((current - previous) / previous) * 100
}

// formatChangeText formats change value as text ("Упало на 5%", "Выросло на 64%")
func (s *MarketingService) formatChangeText(change float64) string {
	if change == 0 {
		return "Без изменений"
	}

	absChange := change
	if change < 0 {
		absChange = -change
	}

	if change > 0 {
		return fmt.Sprintf("Выросло на %.0f%%", absChange)
	}
	return fmt.Sprintf("Упало на %.0f%%", absChange)
}
