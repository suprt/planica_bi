package services

import (
	"fmt"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/integrations"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gorm.io/gorm"
)

// SyncService handles data synchronization with Yandex APIs
type SyncService struct {
	projectRepo  *repositories.ProjectRepository
	metricsRepo  *repositories.MetricsRepository
	directRepo   *repositories.DirectRepository
	metricaClient *integrations.YandexMetricaClient
	directClient  *integrations.YandexDirectClient
	db           *gorm.DB
}

// NewSyncService creates a new sync service
func NewSyncService(
	projectRepo *repositories.ProjectRepository,
	metricsRepo *repositories.MetricsRepository,
	directRepo *repositories.DirectRepository,
	metricaClient *integrations.YandexMetricaClient,
	directClient *integrations.YandexDirectClient,
	db *gorm.DB,
) *SyncService {
	return &SyncService{
		projectRepo:   projectRepo,
		metricsRepo:   metricsRepo,
		directRepo:    directRepo,
		metricaClient: metricaClient,
		directClient:  directClient,
		db:            db,
	}
}

// SyncProject synchronizes data for a specific project
func (s *SyncService) SyncProject(projectID uint) error {
	// Get project
	project, err := s.projectRepo.GetByID(projectID)
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
	if err := s.syncMetricaData(projectID, currentYear, currentMonth); err != nil {
		return fmt.Errorf("failed to sync Metrica data: %w", err)
	}

	// Sync Yandex.Direct data
	if err := s.syncDirectData(projectID, currentYear, currentMonth); err != nil {
		return fmt.Errorf("failed to sync Direct data: %w", err)
	}

	return nil
}

// syncMetricaData synchronizes Yandex.Metrica data for a project
func (s *SyncService) syncMetricaData(projectID uint, year, month int) error {
	// Get all counters for the project
	var counters []models.YandexCounter
	if err := s.db.Where("project_id = ?", projectID).Find(&counters).Error; err != nil {
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
	counterCount := 0

	for _, counter := range counters {
		// Get metrics from API
		metricsData, err := s.metricaClient.GetMetrics(counter.CounterID, dateFrom, dateTo)
		if err != nil {
			// Log error but continue with other counters
			continue
		}

		// TODO: Parse metricsData and aggregate
		// For now, this is a placeholder structure
		_ = metricsData

		// Get age breakdown
		ageData, err := s.metricaClient.GetMetricsByAge(counter.CounterID, dateFrom, dateTo)
		if err != nil {
			continue
		}

		// TODO: Parse ageData and save
		_ = ageData

		counterCount++
	}

	// Get goals for conversions
	var counterIDs []uint
	for _, counter := range counters {
		counterIDs = append(counterIDs, counter.ID)
	}
	
	var goals []models.Goal
	if len(counterIDs) > 0 {
		if err := s.db.Where("counter_id IN ?", counterIDs).Find(&goals).Error; err == nil {
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
						conversionsData, err := s.metricaClient.GetConversions(counter.CounterID, goalIDs, dateFrom, dateTo)
						if err == nil {
							// TODO: Parse conversionsData
							_ = conversionsData
						}
						break
					}
				}
			}
		}
	}

	// Save aggregated monthly metrics
	// TODO: Calculate actual aggregated values from API responses
	monthlyMetrics := &models.MetricsMonthly{
		ProjectID:             projectID,
		Year:                  year,
		Month:                 month,
		Visits:                totalVisits,
		Users:                 totalUsers,
		BounceRate:            totalBounceRate,
		AvgSessionDurationSec: totalDurationSec,
	}

	// Check if record exists
	existing, err := s.metricsRepo.GetMonthlyMetrics(projectID, year, month)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing != nil {
		monthlyMetrics.ID = existing.ID
	}

	return s.metricsRepo.SaveMonthlyMetrics(monthlyMetrics)
}

// syncDirectData synchronizes Yandex.Direct data for a project
func (s *SyncService) syncDirectData(projectID uint, year, month int) error {
	// Get all Direct accounts for the project
	var accounts []models.DirectAccount
	if err := s.db.Where("project_id = ?", projectID).Find(&accounts).Error; err != nil {
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
		// Create Direct client for this account
		// TODO: Get token from config or account settings
		// For now, use shared client
		_ = account // account will be used when implementing per-account tokens
		directClient := s.directClient

		// Get campaign report
		reportData, err := directClient.GetCampaignReport(dateFrom, dateTo)
		if err != nil {
			// Log error but continue with other accounts
			continue
		}

		// TODO: Parse reportData and aggregate
		// For now, this is a placeholder structure
		_ = reportData

		// Get campaigns list
		campaignsData, err := directClient.GetCampaigns()
		if err != nil {
			continue
		}

		// TODO: Parse campaignsData and save individual campaign metrics
		_ = campaignsData
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
	existing, err := s.directRepo.GetTotalsMonthly(projectID, year, month)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing != nil {
		totals.ID = existing.ID
	}

	return s.directRepo.SaveTotalsMonthly(totals)
}

// SyncAllProjects synchronizes data for all active projects
func (s *SyncService) SyncAllProjects() error {
	projects, err := s.projectRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	for _, project := range projects {
		if !project.IsActive {
			continue
		}

		if err := s.SyncProject(project.ID); err != nil {
			// Log error but continue with other projects
			// TODO: Add proper logging
			continue
		}
	}

	return nil
}

// FinalizeMonth finalizes data for the previous month
func (s *SyncService) FinalizeMonth() error {
	// Called on 1st of each month at 07:00 MSK
	// Finalize data for the previous month
	now := time.Now()
	prevMonth := now.AddDate(0, -1, 0)
	year := prevMonth.Year()
	month := int(prevMonth.Month())

	// Get all active projects
	projects, err := s.projectRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	for _, project := range projects {
		if !project.IsActive {
			continue
		}

		// Sync data for previous month to ensure it's finalized
		if err := s.syncMetricaData(project.ID, year, month); err != nil {
			// Log error but continue
			continue
		}

		if err := s.syncDirectData(project.ID, year, month); err != nil {
			// Log error but continue
			continue
		}
	}

	return nil
}

