package cron

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/queue"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"go.uber.org/zap"
)

// Scheduler handles scheduled tasks using cron
type Scheduler struct {
	cron        *cron.Cron
	queueClient *queue.Client
	projectRepo *repositories.ProjectRepository
}

// NewScheduler creates a new scheduler
func NewScheduler(queueClient *queue.Client, projectRepo *repositories.ProjectRepository) *Scheduler {
	// Create cron with timezone support (MSK = Europe/Moscow)
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		// Fallback to UTC if timezone loading fails
		loc = time.UTC
		if logger.Log != nil {
			logger.Log.Warn("Failed to load Europe/Moscow timezone, using UTC", zap.Error(err))
		}
	}

	c := cron.New(cron.WithLocation(loc), cron.WithSeconds())

	return &Scheduler{
		cron:        c,
		queueClient: queueClient,
		projectRepo: projectRepo,
	}
}

// StartDailySync starts daily synchronization task
// Runs every night at 02:00 MSK to update current month (M) data
func (s *Scheduler) StartDailySync() {
	// Schedule: daily at 02:00 MSK (0 0 2 * * *)
	_, err := s.cron.AddFunc("0 0 2 * * *", func() {
		s.runDailySync()
	})
	if err != nil {
		if logger.Log != nil {
			logger.Log.Fatal("Failed to schedule daily sync", zap.Error(err))
		}
		return
	}

	if logger.Log != nil {
		logger.Log.Info("Daily sync scheduled", zap.String("schedule", "02:00 MSK daily"))
	}
}

// StartMonthlyFinalization starts monthly finalization task
// Runs on 1st of each month at 07:00 MSK
func (s *Scheduler) StartMonthlyFinalization() {
	// Schedule: 1st day of month at 07:00 MSK (0 0 7 1 * *)
	_, err := s.cron.AddFunc("0 0 7 1 * *", func() {
		s.runMonthlyFinalization()
	})
	if err != nil {
		if logger.Log != nil {
			logger.Log.Fatal("Failed to schedule monthly finalization", zap.Error(err))
		}
		return
	}

	if logger.Log != nil {
		logger.Log.Info("Monthly finalization scheduled", zap.String("schedule", "07:00 MSK on 1st of month"))
	}
}

// Start starts the cron scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	if logger.Log != nil {
		logger.Log.Info("Cron scheduler started")
	}
}

// Stop stops the cron scheduler
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	if logger.Log != nil {
		logger.Log.Info("Cron scheduler stopped")
	}
}

// runDailySync executes daily synchronization for all active projects
func (s *Scheduler) runDailySync() {
	if logger.Log != nil {
		logger.Log.Info("Starting daily sync for all projects")
	}

	ctx := context.Background()
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("Failed to get projects for daily sync", zap.Error(err))
		}
		return
	}

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	successCount := 0
	errorCount := 0

	for _, project := range projects {
		if !project.IsActive {
			continue
		}

		// Enqueue sync task for each project
		_, err := s.queueClient.EnqueueSyncProjectTask(project.ID)
		if err != nil {
			errorCount++
			if logger.Log != nil {
				logger.Log.Error("Failed to enqueue sync task for project",
					zap.Uint("project_id", project.ID),
					zap.Error(err),
				)
			}
			continue
		}

		successCount++
		if logger.Log != nil {
			logger.Log.Info("Enqueued sync task for project",
				zap.Uint("project_id", project.ID),
				zap.Int("year", year),
				zap.Int("month", month),
			)
		}
	}

	if logger.Log != nil {
		logger.Log.Info("Daily sync completed",
			zap.Int("total_projects", len(projects)),
			zap.Int("success", successCount),
			zap.Int("errors", errorCount),
		)
	}
}

// runMonthlyFinalization executes monthly finalization for all active projects
func (s *Scheduler) runMonthlyFinalization() {
	if logger.Log != nil {
		logger.Log.Info("Starting monthly finalization for all projects")
	}

	ctx := context.Background()
	projects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("Failed to get projects for monthly finalization", zap.Error(err))
		}
		return
	}

	// Calculate previous month
	now := time.Now()
	prevMonth := now.AddDate(0, -1, 0)
	year := prevMonth.Year()
	month := int(prevMonth.Month())

	successCount := 0
	errorCount := 0

	for _, project := range projects {
		if !project.IsActive {
			continue
		}

		// Enqueue sync tasks for previous month
		_, err := s.queueClient.EnqueueSyncMetricaTask(project.ID, year, month)
		if err != nil {
			errorCount++
			if logger.Log != nil {
				logger.Log.Error("Failed to enqueue Metrica finalization task",
					zap.Uint("project_id", project.ID),
					zap.Error(err),
				)
			}
		} else {
			successCount++
		}

		_, err = s.queueClient.EnqueueSyncDirectTask(project.ID, year, month)
		if err != nil {
			errorCount++
			if logger.Log != nil {
				logger.Log.Error("Failed to enqueue Direct finalization task",
					zap.Uint("project_id", project.ID),
					zap.Error(err),
				)
			}
		} else {
			successCount++
		}

		if logger.Log != nil {
			logger.Log.Info("Enqueued finalization tasks for project",
				zap.Uint("project_id", project.ID),
				zap.Int("year", year),
				zap.Int("month", month),
			)
		}
	}

	if logger.Log != nil {
		logger.Log.Info("Monthly finalization completed",
			zap.Int("total_projects", len(projects)),
			zap.Int("success", successCount),
			zap.Int("errors", errorCount),
			zap.Int("year", year),
			zap.Int("month", month),
		)
	}
}
