package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ReportService handles business logic for reports
type ReportService struct {
	metricsRepo *repositories.MetricsRepository
	directRepo  *repositories.DirectRepository
	seoRepo     *repositories.SEORepository
	projectRepo *repositories.ProjectRepository
	cfg         *config.Config
}

// NewReportService creates a new report service
func NewReportService(
	metricsRepo *repositories.MetricsRepository,
	directRepo *repositories.DirectRepository,
	seoRepo *repositories.SEORepository,
	projectRepo *repositories.ProjectRepository,
	cfg *config.Config,
) *ReportService {
	return &ReportService{
		metricsRepo: metricsRepo,
		directRepo:  directRepo,
		seoRepo:     seoRepo,
		projectRepo: projectRepo,
		cfg:         cfg,
	}
}

// Dynamics represents percentage change compared to previous period
type Dynamics struct {
	Visits float64 `json:"visits"`
	Users  float64 `json:"users"`
	Bounce float64 `json:"bounce"`
	AvgSec float64 `json:"avgSec"`
	Conv   float64 `json:"conv,omitempty"`
}

// MetricaSummaryRow represents a single row in metrica summary
type MetricaSummaryRow struct {
	Month    string    `json:"month"`
	Visits   int       `json:"visits"`
	Users    int       `json:"users"`
	Bounce   float64   `json:"bounce"`
	AvgSec   int       `json:"avgSec"`
	Conv     *int      `json:"conv,omitempty"`
	Dynamics *Dynamics `json:"dynamics,omitempty"`
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

// AiInsights represents AI analysis insights
type AiInsights struct {
	Summary         string   `json:"summary"`
	Recommendations []string `json:"recommendations,omitempty"`
}

// Report represents a full report according to TZ format
type Report struct {
	ProjectID  uint        `json:"projectId"`
	Periods    []string    `json:"periods"`
	Metrica    MetricaData `json:"metrica"`
	Direct     DirectData  `json:"direct"`
	SEO        SEOData     `json:"seo"`
	AiInsights *AiInsights `json:"ai_insights,omitempty"`
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

		// Get SEO summary (organic visitors and conversions)
		// For now, we can use Metrica metrics as a base, but ideally should filter by organic traffic
		// TODO: Add proper organic traffic filtering via Metrica API with source segment
		metrics, err := s.metricsRepo.GetMonthlyMetrics(ctx, projectID, pd.year, pd.month)
		if err == nil && metrics != nil {
			// Calculate organic visitors as a percentage of total (rough estimate: 30-50% is typical for organic)
			// In production, this should come from Metrica API with organic segment filter
			organicVisitors := int(float64(metrics.Users) * 0.4) // Rough estimate: 40% organic
			organicConversions := 0
			if metrics.Conversions != nil {
				organicConversions = int(float64(*metrics.Conversions) * 0.4) // Same percentage for conversions
			}

			report.SEO.Summary = append(report.SEO.Summary, SEOSummaryRow{
				Month:    pd.period,
				Visitors: organicVisitors,
				Conv:     organicConversions,
			})
		}
	}

	// Convert campaign map to slice
	for _, campaignData := range campaignMap {
		report.Direct.Campaigns = append(report.Direct.Campaigns, *campaignData)
	}

	// Calculate dynamics for Metrica summary (compare current month with previous)
	if len(report.Metrica.Summary) >= 2 {
		m0 := &report.Metrica.Summary[0] // Current month
		m1 := &report.Metrica.Summary[1] // Previous month

		m0.Dynamics = &Dynamics{
			Visits: utils.CalculateDynamics(float64(m0.Visits), float64(m1.Visits)),
			Users:  utils.CalculateDynamics(float64(m0.Users), float64(m1.Users)),
			Bounce: utils.CalculateDynamics(m0.Bounce, m1.Bounce),
			AvgSec: utils.CalculateDynamics(float64(m0.AvgSec), float64(m1.AvgSec)),
		}

		// Calculate dynamics for conversions if both have values
		var conv0, conv1 float64
		if m0.Conv != nil {
			conv0 = float64(*m0.Conv)
		}
		if m1.Conv != nil {
			conv1 = float64(*m1.Conv)
		}
		if conv1 > 0 {
			m0.Dynamics.Conv = utils.CalculateDynamics(conv0, conv1)
		}
	}

	return report, nil
}

// CalculateDynamics calculates percentage change between two values
func (s *ReportService) CalculateDynamics(current, previous float64) float64 {
	return utils.CalculateDynamics(current, previous)
}

// ChannelMetrics represents metrics for a channel
type ChannelMetrics struct {
	CPC         []float64 `json:"cpc"`
	Impressions []int     `json:"impressions"`
	Clicks      []int     `json:"clicks"`
	CTR         []float64 `json:"ctr"`
	Conversions []int     `json:"conversions"`
	CPA         []float64 `json:"cpa"`
	Cost        []float64 `json:"cost"` // 💰 бюджет в рублях
}

// ChannelMetricsOutput represents the output format for channel metrics
type ChannelMetricsOutput struct {
	Project string                     `json:"project"`
	Periods []string                   `json:"periods"`
	Metrics map[string]*ChannelMetrics `json:"metrics"`
}

// GetChannelMetrics retrieves channel metrics from database for specified periods
func (s *ReportService) GetChannelMetrics(ctx context.Context, projectID uint, periods []string) (*ChannelMetricsOutput, error) {
	// Get project to get project name
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	output := &ChannelMetricsOutput{
		Project: project.Name,
		Periods: periods,
		Metrics: make(map[string]*ChannelMetrics),
	}

	// Initialize channel metrics
	simpleMetrics := &ChannelMetrics{}
	mkMetrics := &ChannelMetrics{}
	rsyaMetrics := &ChannelMetrics{}

	// Process each period
	for _, period := range periods {
		year, month, err := parsePeriod(period)
		if err != nil {
			return nil, fmt.Errorf("invalid period format %s: %w", period, err)
		}

		// Get "simple" channel data (Direct totals)
		directTotals, err := s.directRepo.GetTotalsMonthly(ctx, projectID, year, month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to get direct totals: %w", err)
		}
		if directTotals != nil {
			simpleMetrics.CPC = append(simpleMetrics.CPC, directTotals.CPC)
			simpleMetrics.Impressions = append(simpleMetrics.Impressions, directTotals.Impressions)
			simpleMetrics.Clicks = append(simpleMetrics.Clicks, directTotals.Clicks)
			simpleMetrics.CTR = append(simpleMetrics.CTR, directTotals.CTRPct)
			if directTotals.Conversions != nil {
				simpleMetrics.Conversions = append(simpleMetrics.Conversions, *directTotals.Conversions)
			} else {
				simpleMetrics.Conversions = append(simpleMetrics.Conversions, 0)
			}
			if directTotals.CPA != nil {
				simpleMetrics.CPA = append(simpleMetrics.CPA, *directTotals.CPA)
			} else {
				simpleMetrics.CPA = append(simpleMetrics.CPA, 0)
			}
			simpleMetrics.Cost = append(simpleMetrics.Cost, directTotals.Cost)
		} else {
			// Add zeros if no data
			simpleMetrics.CPC = append(simpleMetrics.CPC, 0)
			simpleMetrics.Impressions = append(simpleMetrics.Impressions, 0)
			simpleMetrics.Clicks = append(simpleMetrics.Clicks, 0)
			simpleMetrics.CTR = append(simpleMetrics.CTR, 0)
			simpleMetrics.Conversions = append(simpleMetrics.Conversions, 0)
			simpleMetrics.CPA = append(simpleMetrics.CPA, 0)
			simpleMetrics.Cost = append(simpleMetrics.Cost, 0)
		}

		// Get "МК" channel data (Metrica metrics)
		metrics, err := s.metricsRepo.GetMonthlyMetrics(ctx, projectID, year, month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to get metrica metrics: %w", err)
		}
		if metrics != nil {
			// For Metrica, we don't have CPC/CPA directly, so we'll use 0 or calculate from other data
			// Metrica doesn't have impressions/clicks in the same way, so we'll use visits/users
			mkMetrics.CPC = append(mkMetrics.CPC, 0)                              // Metrica doesn't have CPC
			mkMetrics.Impressions = append(mkMetrics.Impressions, metrics.Visits) // Using visits as approximation
			mkMetrics.Clicks = append(mkMetrics.Clicks, metrics.Users)            // Using users as approximation
			mkMetrics.CTR = append(mkMetrics.CTR, 0)                              // Metrica doesn't have CTR in this context
			if metrics.Conversions != nil {
				mkMetrics.Conversions = append(mkMetrics.Conversions, *metrics.Conversions)
			} else {
				mkMetrics.Conversions = append(mkMetrics.Conversions, 0)
			}
			mkMetrics.CPA = append(mkMetrics.CPA, 0)   // Metrica doesn't have CPA
			mkMetrics.Cost = append(mkMetrics.Cost, 0) // Metrica doesn't track cost
		} else {
			// Add zeros if no data
			mkMetrics.CPC = append(mkMetrics.CPC, 0)
			mkMetrics.Impressions = append(mkMetrics.Impressions, 0)
			mkMetrics.Clicks = append(mkMetrics.Clicks, 0)
			mkMetrics.CTR = append(mkMetrics.CTR, 0)
			mkMetrics.Conversions = append(mkMetrics.Conversions, 0)
			mkMetrics.CPA = append(mkMetrics.CPA, 0)
			mkMetrics.Cost = append(mkMetrics.Cost, 0)
		}

		// Get "РСЯ" channel data (sum of all Direct campaigns)
		directCampaigns, err := s.directRepo.GetCampaignMonthly(ctx, projectID, year, month)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to get direct campaigns: %w", err)
		}

		var rsyaImpressions, rsyaClicks, rsyaConversions int
		var rsyaCost, rsyaCPC, rsyaCTR, rsyaCPA float64

		for _, campaign := range directCampaigns {
			rsyaImpressions += campaign.Impressions
			rsyaClicks += campaign.Clicks
			rsyaCost += campaign.Cost
			if campaign.Conversions != nil {
				rsyaConversions += *campaign.Conversions
			}
		}

		// Calculate aggregated metrics for РСЯ
		if rsyaClicks > 0 {
			rsyaCTR = (float64(rsyaClicks) / float64(rsyaImpressions)) * 100
			rsyaCPC = rsyaCost / float64(rsyaClicks)
		}
		if rsyaConversions > 0 {
			rsyaCPA = rsyaCost / float64(rsyaConversions)
		}

		rsyaMetrics.CPC = append(rsyaMetrics.CPC, rsyaCPC)
		rsyaMetrics.Impressions = append(rsyaMetrics.Impressions, rsyaImpressions)
		rsyaMetrics.Clicks = append(rsyaMetrics.Clicks, rsyaClicks)
		rsyaMetrics.CTR = append(rsyaMetrics.CTR, rsyaCTR)
		rsyaMetrics.Conversions = append(rsyaMetrics.Conversions, rsyaConversions)
		rsyaMetrics.CPA = append(rsyaMetrics.CPA, rsyaCPA)
		rsyaMetrics.Cost = append(rsyaMetrics.Cost, rsyaCost)
	}

	// Add metrics to output
	output.Metrics["simple"] = simpleMetrics
	output.Metrics["МК"] = mkMetrics
	output.Metrics["РСЯ"] = rsyaMetrics

	return output, nil
}

// parsePeriod parses a period string "YYYY-MM" into year and month
func parsePeriod(period string) (year int, month int, error error) {
	parts := strings.Split(period, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid period format, expected YYYY-MM")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year: %w", err)
	}

	month, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid month: %w", err)
	}

	if month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("month must be between 1 and 12")
	}

	return year, month, nil
}

// MetricsAnalysisResult represents the result of metrics analysis
type MetricsAnalysisResult struct {
	AnalyticalFacts string `json:"analytical_facts"`
	AIReport        string `json:"ai_report,omitempty"`
	Error           string `json:"error,omitempty"`
}

// AnalyzeChannelMetrics analyzes channel metrics using Python script and OpenAI
func (s *ReportService) AnalyzeChannelMetrics(ctx context.Context, metricsData *ChannelMetricsOutput) (*MetricsAnalysisResult, error) {
	// Convert metrics data to JSON
	jsonData, err := json.Marshal(metricsData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metrics data: %w", err)
	}

	// Find Python script path
	// Try multiple possible locations
	possiblePaths := []string{
		filepath.Join("scripts", "analyze_metrics.py"),
		filepath.Join("backend", "scripts", "analyze_metrics.py"),
		filepath.Join("/root", "scripts", "analyze_metrics.py"), // Docker container path
		"./scripts/analyze_metrics.py",
	}

	var scriptPath string
	var found bool
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			scriptPath = path
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("python script not found. Tried: %v", possiblePaths)
	}

	// Create context with timeout for script execution (120 seconds)
	scriptCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Create command to run Python script
	cmd := exec.CommandContext(scriptCtx, "python3", scriptPath)
	// If python3 is not available, try python
	if _, err := exec.LookPath("python3"); err != nil {
		cmd = exec.CommandContext(scriptCtx, "python", scriptPath)
	}

	// Set environment variables for Ollama API
	cmd.Env = os.Environ()
	if s.cfg.OllamaAPIKey != "" {
		cmd.Env = append(cmd.Env, "OLLAMA_API_KEY="+s.cfg.OllamaAPIKey)
	}
	if s.cfg.OllamaAPIURL != "" {
		cmd.Env = append(cmd.Env, "OLLAMA_API_URL="+s.cfg.OllamaAPIURL)
	}
	if s.cfg.OllamaModel != "" {
		cmd.Env = append(cmd.Env, "OLLAMA_MODEL="+s.cfg.OllamaModel)
	}

	// Set stdin to JSON data
	cmd.Stdin = bytes.NewReader(jsonData)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the script
	if err := cmd.Run(); err != nil {
		if logger.Log != nil {
			logger.Log.Error("Failed to run Python analysis script",
				zap.Error(err),
				zap.String("stderr", stderr.String()),
			)
		}
		return nil, fmt.Errorf("python script execution failed: %w (stderr: %s)", err, stderr.String())
	}

	// Parse JSON output
	var result MetricsAnalysisResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Python script output: %w (output: %s)", err, stdout.String())
	}

	return &result, nil
}
