package services

import (
	"context"
	"fmt"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/integrations"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SyncService handles data synchronization with Yandex APIs
type SyncService struct {
	projectRepo   *repositories.ProjectRepository
	metricsRepo   *repositories.MetricsRepository
	directRepo    *repositories.DirectRepository
	counterRepo   *repositories.CounterRepository
	goalRepo      *repositories.GoalRepository
	metricaClient *integrations.YandexMetricaClient
	directClient  *integrations.YandexDirectClient
}

// NewSyncService creates a new sync service
func NewSyncService(
	projectRepo *repositories.ProjectRepository,
	metricsRepo *repositories.MetricsRepository,
	directRepo *repositories.DirectRepository,
	counterRepo *repositories.CounterRepository,
	goalRepo *repositories.GoalRepository,
	metricaClient *integrations.YandexMetricaClient,
	directClient *integrations.YandexDirectClient,
) *SyncService {
	return &SyncService{
		projectRepo:   projectRepo,
		metricsRepo:   metricsRepo,
		directRepo:    directRepo,
		counterRepo:   counterRepo,
		goalRepo:      goalRepo,
		metricaClient: metricaClient,
		directClient:  directClient,
	}
}

// SyncProject synchronizes data for a specific project
func (s *SyncService) SyncProject(ctx context.Context, projectID uint) error {
	// Get project
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	if !project.IsActive {
		return fmt.Errorf("project %d is not active", projectID)
	}

	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// Sync Yandex.Metrica data
	if err := s.syncMetricaData(ctx, projectID, currentYear, currentMonth); err != nil {
		return fmt.Errorf("failed to sync Metrica data: %w", err)
	}

	// Sync Yandex.Direct data
	if err := s.syncDirectData(ctx, projectID, currentYear, currentMonth); err != nil {
		return fmt.Errorf("failed to sync Direct data: %w", err)
	}

	return nil
}

// syncMetricaData synchronizes Yandex.Metrica data for a project
func (s *SyncService) syncMetricaData(ctx context.Context, projectID uint, year, month int) error {
	// Get all counters for the project
	counters, err := s.counterRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get counters: %w", err)
	}

	if len(counters) == 0 {
		return nil // No counters to sync
	}

	// Calculate date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).AddDate(0, 0, -1) // Last day of month
	dateFrom := startDate.Format("2006-01-02")
	dateTo := endDate.Format("2006-01-02")

	// Aggregate metrics from all counters
	var totalVisits, totalUsers int
	var totalBounceRate float64
	var totalDurationSec int
	var bounceRateSum float64
	counterCount := 0

	for _, counter := range counters {
		// Get metrics from API
		metricsData, err := s.metricaClient.GetMetrics(ctx, counter.CounterID, dateFrom, dateTo)
		if err != nil {
			// Log error but continue with other counters
			if logger.Log != nil {
				logger.Log.Warn("Failed to get metrics from Metrica API",
					zap.Int64("counter_id", counter.CounterID),
					zap.Error(err),
				)
			}
			continue
		}

		// Parse metricsData and aggregate
		visits, users, bounceRate, duration := s.parseMetricaMetrics(metricsData)
		totalVisits += visits
		totalUsers += users
		if bounceRate > 0 {
			bounceRateSum += bounceRate
			counterCount++
		}
		totalDurationSec += duration

		// Get age breakdown
		ageData, err := s.metricaClient.GetMetricsByAge(ctx, counter.CounterID, dateFrom, dateTo)
		if err != nil {
			if logger.Log != nil {
				logger.Log.Warn("Failed to get age metrics from Metrica API",
					zap.Int64("counter_id", counter.CounterID),
					zap.Error(err),
				)
			}
			continue
		}

		// Parse ageData and save
		if err := s.parseAndSaveAgeMetrics(ctx, ageData, projectID, counter.ID, year, month); err != nil {
			if logger.Log != nil {
				logger.Log.Warn("Failed to save age metrics",
					zap.Int64("counter_id", counter.CounterID),
					zap.Error(err),
				)
			}
		}
	}

	// Calculate average bounce rate
	if counterCount > 0 {
		totalBounceRate = bounceRateSum / float64(counterCount)
	}

	// Initialize monthly metrics
	monthlyMetrics := &models.MetricsMonthly{
		ProjectID:             projectID,
		Year:                  year,
		Month:                 month,
		Visits:                totalVisits,
		Users:                 totalUsers,
		BounceRate:            totalBounceRate,
		AvgSessionDurationSec: totalDurationSec,
	}

	// Get goals for conversions
	var counterIDs []uint
	for _, counter := range counters {
		counterIDs = append(counterIDs, counter.ID)
	}

	if len(counterIDs) > 0 {
		goals, err := s.goalRepo.GetByCounterIDs(ctx, counterIDs)
		if err == nil {
			goalIDs := make([]int64, 0, len(goals))
			for _, goal := range goals {
				if goal.IsConversion {
					goalIDs = append(goalIDs, goal.GoalID)
				}
			}

			if len(goalIDs) > 0 {
				// Get conversions for primary counter
				for _, counter := range counters {
					if counter.IsPrimary {
						conversionsData, err := s.metricaClient.GetConversions(ctx, counter.CounterID, goalIDs, dateFrom, dateTo)
						if err == nil {
							conversions := s.parseConversions(conversionsData)
							if conversions != nil {
								monthlyMetrics.Conversions = conversions
							}
						} else if logger.Log != nil {
							logger.Log.Warn("Failed to get conversions from Metrica API",
								zap.Int64("counter_id", counter.CounterID),
								zap.Error(err),
							)
						}
						break
					}
				}
			}
		}
	}

	// Check if record exists
	existing, err := s.metricsRepo.GetMonthlyMetrics(ctx, projectID, year, month)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing != nil {
		monthlyMetrics.ID = existing.ID
	}

	return s.metricsRepo.SaveMonthlyMetrics(monthlyMetrics)
}

// syncDirectData synchronizes Yandex.Direct data for a project
func (s *SyncService) syncDirectData(ctx context.Context, projectID uint, year, month int) error {
	// Get all Direct accounts for the project
	accounts, err := s.directRepo.GetAccountsByProjectID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get Direct accounts: %w", err)
	}

	if len(accounts) == 0 {
		return nil // No accounts to sync
	}

	// Calculate date range for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).AddDate(0, 0, -1) // Last day of month
	dateFrom := startDate.Format("2006-01-02")
	dateTo := endDate.Format("2006-01-02")

	// Aggregate totals from all accounts
	var totalImpressions, totalClicks int
	var totalCost float64
	var totalCTR float64
	var totalCPC float64

	for _, account := range accounts {
		// Use shared client for now
		// Note: In future, token should be retrieved from account settings or config
		directClient := s.directClient

		// Get campaign report
		reportData, err := directClient.GetCampaignReport(ctx, dateFrom, dateTo)
		if err != nil {
			// Log error but continue with other accounts
			if logger.Log != nil {
				logger.Log.Warn("Failed to get campaign report from Direct API",
					zap.Uint("account_id", account.ID),
					zap.String("client_login", account.ClientLogin),
					zap.Error(err),
				)
			}
			continue
		}

		// Parse reportData and aggregate
		impressions, clicks, cost := s.parseDirectReport(reportData)
		totalImpressions += impressions
		totalClicks += clicks
		totalCost += cost

		// Get campaigns list
		campaignsData, err := directClient.GetCampaigns(ctx)
		if err != nil {
			if logger.Log != nil {
				logger.Log.Warn("Failed to get campaigns from Direct API",
					zap.Uint("account_id", account.ID),
					zap.Error(err),
				)
			}
			continue
		}

		// Parse campaignsData and save individual campaign metrics
		if err := s.parseAndSaveCampaignMetrics(ctx, campaignsData, account.ID, projectID, year, month); err != nil {
			if logger.Log != nil {
				logger.Log.Warn("Failed to save campaign metrics",
					zap.Uint("account_id", account.ID),
					zap.Error(err),
				)
			}
		}
	}

	// Calculate aggregated metrics
	if totalClicks > 0 {
		totalCTR = (float64(totalClicks) / float64(totalImpressions)) * 100
		totalCPC = totalCost / float64(totalClicks)
	}

	// Save aggregated monthly totals
	totals := &models.DirectTotalsMonthly{
		ProjectID:   projectID,
		Year:        year,
		Month:       month,
		Impressions: totalImpressions,
		Clicks:      totalClicks,
		CTRPct:      totalCTR,
		CPC:         totalCPC,
		Cost:        totalCost,
	}

	// Check if record exists
	existing, err := s.directRepo.GetTotalsMonthly(ctx, projectID, year, month)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing != nil {
		totals.ID = existing.ID
	}

	return s.directRepo.SaveTotalsMonthly(totals)
}

// SyncAllProjects synchronizes data for all active projects
func (s *SyncService) SyncAllProjects(ctx context.Context) error {
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	for _, project := range projects {
		if !project.IsActive {
			continue
		}

		if err := s.SyncProject(ctx, project.ID); err != nil {
			// Log error but continue with other projects
			if logger.Log != nil {
				logger.Log.Error("Failed to sync project",
					zap.Uint("project_id", project.ID),
					zap.Error(err),
				)
			}
			continue
		}
	}

	return nil
}

// FinalizeMonth finalizes data for the previous month
func (s *SyncService) FinalizeMonth(ctx context.Context) error {
	// Called on 1st of each month at 07:00 MSK
	// Finalize data for the previous month
	now := time.Now()
	prevMonth := now.AddDate(0, -1, 0)
	year := prevMonth.Year()
	month := int(prevMonth.Month())

	// Get all active projects
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	for _, project := range projects {
		if !project.IsActive {
			continue
		}

		// Sync data for previous month to ensure it's finalized
		if err := s.syncMetricaData(ctx, project.ID, year, month); err != nil {
			// Log error but continue
			if logger.Log != nil {
				logger.Log.Error("Failed to finalize Metrica data",
					zap.Uint("project_id", project.ID),
					zap.Int("year", year),
					zap.Int("month", month),
					zap.Error(err),
				)
			}
			continue
		}

		if err := s.syncDirectData(ctx, project.ID, year, month); err != nil {
			// Log error but continue
			if logger.Log != nil {
				logger.Log.Error("Failed to finalize Direct data",
					zap.Uint("project_id", project.ID),
					zap.Int("year", year),
					zap.Int("month", month),
					zap.Error(err),
				)
			}
			continue
		}
	}

	return nil
}

// parseMetricaMetrics parses metrics data from Yandex.Metrica API
// Returns: visits, users, bounceRate, avgDurationSec
func (s *SyncService) parseMetricaMetrics(data interface{}) (int, int, float64, int) {
	if data == nil {
		return 0, 0, 0, 0
	}

	// Try to parse as map[string]interface{}
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return 0, 0, 0, 0
	}

	var visits, users int
	var bounceRate float64
	var durationSec int

	// Parse data array if exists
	if dataArray, ok := dataMap["data"].([]interface{}); ok && len(dataArray) > 0 {
		if firstRow, ok := dataArray[0].([]interface{}); ok && len(firstRow) >= 4 {
			// Typical structure: [visits, users, bounceRate, avgDuration]
			if v, ok := firstRow[0].(float64); ok {
				visits = int(v)
			}
			if v, ok := firstRow[1].(float64); ok {
				users = int(v)
			}
			if v, ok := firstRow[2].(float64); ok {
				bounceRate = v
			}
			if v, ok := firstRow[3].(float64); ok {
				durationSec = int(v)
			}
		}
	}

	// Alternative: try direct fields
	if visits == 0 {
		if v, ok := dataMap["visits"].(float64); ok {
			visits = int(v)
		}
	}
	if users == 0 {
		if v, ok := dataMap["users"].(float64); ok {
			users = int(v)
		}
	}
	if bounceRate == 0 {
		if v, ok := dataMap["bounceRate"].(float64); ok {
			bounceRate = v
		}
	}
	if durationSec == 0 {
		if v, ok := dataMap["avgVisitDurationSeconds"].(float64); ok {
			durationSec = int(v)
		}
	}

	return visits, users, bounceRate, durationSec
}

// parseAndSaveAgeMetrics parses age breakdown data and saves to database
func (s *SyncService) parseAndSaveAgeMetrics(ctx context.Context, data interface{}, projectID uint, counterID uint, year, month int) error {
	if data == nil {
		return nil
	}

	// Try to cast to []AgeMetricsResult first (new format)
	if ageResults, ok := data.([]integrations.AgeMetricsResult); ok {
		return s.parseAndSaveAgeMetricsTyped(ctx, ageResults, projectID, counterID, year, month)
	}

	// Fallback to old format (map[string]interface{}) for backward compatibility
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	// Parse data array
	dataArray, ok := dataMap["data"].([]interface{})
	if !ok {
		return nil
	}

	// Map age intervals to our AgeGroup enum
	ageGroupMap := map[string]models.AgeGroup{
		"18-24":   models.AgeGroup1824,
		"25-34":   models.AgeGroup2534,
		"35-44":   models.AgeGroup3544,
		"45-54":   models.AgeGroup4554,
		"55+":     models.AgeGroup55Plus,
		"unknown": models.AgeGroupUnknown,
	}

	for _, row := range dataArray {
		rowArray, ok := row.([]interface{})
		if !ok || len(rowArray) < 5 {
			continue
		}

		// Parse age group
		ageInterval, ok := rowArray[0].(string)
		if !ok {
			continue
		}
		ageGroup, exists := ageGroupMap[ageInterval]
		if !exists {
			ageGroup = models.AgeGroupUnknown
		}

		// Parse metrics
		var visits, users int
		var bounceRate float64
		var durationSec int

		if v, ok := rowArray[1].(float64); ok {
			visits = int(v)
		}
		if v, ok := rowArray[2].(float64); ok {
			users = int(v)
		}
		if v, ok := rowArray[3].(float64); ok {
			bounceRate = v
		}
		if v, ok := rowArray[4].(float64); ok {
			durationSec = int(v)
		}

		// Check if record exists
		existing, err := s.metricsRepo.GetAgeMetricsByGroup(ctx, projectID, year, month, string(ageGroup))

		ageMetrics := &models.MetricsAgeMonthly{
			ProjectID:             projectID,
			Year:                  year,
			Month:                 month,
			AgeGroup:              ageGroup,
			Visits:                visits,
			Users:                 users,
			BounceRate:            bounceRate,
			AvgSessionDurationSec: durationSec,
		}

		if err == nil && existing != nil {
			ageMetrics.ID = existing.ID
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err := s.metricsRepo.SaveAgeMetrics(ageMetrics); err != nil {
			return err
		}
	}

	return nil
}

// parseAndSaveAgeMetricsTyped parses typed age breakdown data and saves to database
func (s *SyncService) parseAndSaveAgeMetricsTyped(ctx context.Context, ageResults []integrations.AgeMetricsResult, projectID uint, counterID uint, year, month int) error {
	// Map age intervals to our AgeGroup enum
	ageGroupMap := map[string]models.AgeGroup{
		"18-24":   models.AgeGroup1824,
		"25-34":   models.AgeGroup2534,
		"35-44":   models.AgeGroup3544,
		"45-54":   models.AgeGroup4554,
		"55+":     models.AgeGroup55Plus,
		"unknown": models.AgeGroupUnknown,
	}

	for _, result := range ageResults {
		ageGroup, exists := ageGroupMap[result.AgeGroup]
		if !exists {
			ageGroup = models.AgeGroupUnknown
		}

		// Check if record exists
		existing, err := s.metricsRepo.GetAgeMetricsByGroup(ctx, projectID, year, month, string(ageGroup))

		ageMetrics := &models.MetricsAgeMonthly{
			ProjectID:             projectID,
			Year:                  year,
			Month:                 month,
			AgeGroup:              ageGroup,
			Visits:                int(result.Visits),
			Users:                 int(result.Users),
			BounceRate:            result.BounceRate,
			AvgSessionDurationSec: result.AvgSessionDurationSec,
		}

		if err == nil && existing != nil {
			ageMetrics.ID = existing.ID
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err := s.metricsRepo.SaveAgeMetrics(ageMetrics); err != nil {
			return err
		}
	}

	return nil
}

// parseConversions parses conversions data from Yandex.Metrica API
func (s *SyncService) parseConversions(data interface{}) *int {
	if data == nil {
		return nil
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	var totalConversions int

	// Parse data array
	if dataArray, ok := dataMap["data"].([]interface{}); ok {
		for _, row := range dataArray {
			rowArray, ok := row.([]interface{})
			if !ok || len(rowArray) < 2 {
				continue
			}
			// Sum all conversions
			if v, ok := rowArray[1].(float64); ok {
				totalConversions += int(v)
			}
		}
	}

	// Alternative: direct field
	if totalConversions == 0 {
		if v, ok := dataMap["conversions"].(float64); ok {
			totalConversions = int(v)
		}
	}

	if totalConversions == 0 {
		return nil
	}

	return &totalConversions
}

// parseDirectReport parses campaign report data from Yandex.Direct API
// Returns: impressions, clicks, cost
func (s *SyncService) parseDirectReport(data interface{}) (int, int, float64) {
	if data == nil {
		return 0, 0, 0
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return 0, 0, 0
	}

	var impressions, clicks int
	var cost float64

	// Parse report structure
	if report, ok := dataMap["report"].(map[string]interface{}); ok {
		if rows, ok := report["rows"].([]interface{}); ok {
			for _, row := range rows {
				rowMap, ok := row.(map[string]interface{})
				if !ok {
					continue
				}

				if v, ok := rowMap["Impressions"].(float64); ok {
					impressions += int(v)
				}
				if v, ok := rowMap["Clicks"].(float64); ok {
					clicks += int(v)
				}
				if v, ok := rowMap["Cost"].(float64); ok {
					cost += v
				}
			}
		}
	}

	// Alternative: direct totals
	if impressions == 0 {
		if v, ok := dataMap["Impressions"].(float64); ok {
			impressions = int(v)
		}
	}
	if clicks == 0 {
		if v, ok := dataMap["Clicks"].(float64); ok {
			clicks = int(v)
		}
	}
	if cost == 0 {
		if v, ok := dataMap["Cost"].(float64); ok {
			cost = v
		}
	}

	return impressions, clicks, cost
}

// parseAndSaveCampaignMetrics parses campaigns data and saves individual campaign metrics
func (s *SyncService) parseAndSaveCampaignMetrics(ctx context.Context, data interface{}, accountID uint, projectID uint, year, month int) error {
	if data == nil {
		return nil
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	// Get campaigns array
	var campaigns []interface{}
	if result, ok := dataMap["result"].(map[string]interface{}); ok {
		if campaignsArray, ok := result["Campaigns"].([]interface{}); ok {
			campaigns = campaignsArray
		}
	} else if campaignsArray, ok := dataMap["Campaigns"].([]interface{}); ok {
		campaigns = campaignsArray
	}

	if len(campaigns) == 0 {
		return nil
	}

	// Get DirectCampaign records for this account
	dbCampaigns, err := s.directRepo.GetCampaignsByAccountID(ctx, accountID)
	if err != nil {
		return err
	}

	// Create map for quick lookup
	campaignMap := make(map[int64]uint)
	for _, dbCampaign := range dbCampaigns {
		campaignMap[dbCampaign.CampaignID] = dbCampaign.ID
	}

	// Parse and save metrics for each campaign
	for _, campaignData := range campaigns {
		campaignMapData, ok := campaignData.(map[string]interface{})
		if !ok {
			continue
		}

		var campaignID int64
		if id, ok := campaignMapData["Id"].(float64); ok {
			campaignID = int64(id)
		} else if id, ok := campaignMapData["CampaignId"].(float64); ok {
			campaignID = int64(id)
		}

		if campaignID == 0 {
			continue
		}

		// Find or create DirectCampaign record
		var directCampaignID uint
		if dbID, exists := campaignMap[campaignID]; exists {
			directCampaignID = dbID
		} else {
			// Create new campaign record
			newCampaign := &models.DirectCampaign{
				DirectAccountID: accountID,
				CampaignID:      campaignID,
			}
			if name, ok := campaignMapData["Name"].(string); ok {
				newCampaign.Name = name
			}
			if status, ok := campaignMapData["Status"].(string); ok {
				newCampaign.Status = status
			}
			if err := s.directRepo.CreateCampaign(ctx, newCampaign); err != nil {
				continue
			}
			directCampaignID = newCampaign.ID
		}

		// Parse metrics (if available in campaign data)
		var impressions, clicks int
		var cost float64
		var ctr, cpc float64

		if stats, ok := campaignMapData["Statistics"].(map[string]interface{}); ok {
			if v, ok := stats["Impressions"].(float64); ok {
				impressions = int(v)
			}
			if v, ok := stats["Clicks"].(float64); ok {
				clicks = int(v)
			}
			if v, ok := stats["Cost"].(float64); ok {
				cost = v
			}
		}

		// Calculate CTR and CPC
		if impressions > 0 && clicks > 0 {
			ctr = (float64(clicks) / float64(impressions)) * 100
		}
		if clicks > 0 && cost > 0 {
			cpc = cost / float64(clicks)
		}

		// Save campaign monthly metrics
		existing, err := s.directRepo.GetCampaignMonthlyByCampaignID(ctx, projectID, directCampaignID, year, month)

		campaignMetrics := &models.DirectCampaignMonthly{
			ProjectID:        projectID,
			DirectCampaignID: directCampaignID,
			Year:             year,
			Month:            month,
			Impressions:      impressions,
			Clicks:           clicks,
			CTRPct:           ctr,
			CPC:              cpc,
			Cost:             cost,
		}

		if err == nil && existing != nil {
			campaignMetrics.ID = existing.ID
		} else if err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		if err := s.directRepo.SaveCampaignMonthly(campaignMetrics); err != nil {
			continue
		}
	}

	return nil
}
