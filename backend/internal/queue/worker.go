package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/cache"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/services"
	"go.uber.org/zap"
)

// Worker handles task processing
type Worker struct {
	server        *asynq.Server
	mux           *asynq.ServeMux
	syncService   *services.SyncService
	reportService *services.ReportService
	cache         *cache.Cache
}

// NewWorker creates a new queue worker
func NewWorker(cfg *config.Config, syncService *services.SyncService, reportService *services.ReportService, cacheClient *cache.Cache) (*Worker, error) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}

	// Create server with config
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10, // Process up to 10 tasks concurrently
			Queues: map[string]int{
				"critical": 6, // High priority
				"default":  3, // Normal priority
				"low":      1, // Low priority
			},
		},
	)

	// Create mux for routing tasks
	mux := asynq.NewServeMux()

	worker := &Worker{
		server:        server,
		mux:           mux,
		syncService:   syncService,
		reportService: reportService,
		cache:         cacheClient,
	}

	// Register task handlers
	worker.registerHandlers()

	return worker, nil
}

// registerHandlers registers all task handlers
func (w *Worker) registerHandlers() {
	w.mux.HandleFunc(TypeSyncMetrica, w.handleSyncMetrica)
	w.mux.HandleFunc(TypeSyncDirect, w.handleSyncDirect)
	w.mux.HandleFunc(TypeSyncProject, w.handleSyncProject)
	w.mux.HandleFunc(TypeAnalyzeMetrics, w.handleAnalyzeMetrics)
	w.mux.HandleFunc(TypeGenerateReport, w.handleGenerateReport)
}

// handleSyncMetrica handles Metrica sync task
func (w *Worker) handleSyncMetrica(ctx context.Context, task *asynq.Task) error {
	payload, err := ParseSyncMetricaPayload(task)
	if err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	if logger.Log != nil {
		logger.Log.Info("Processing Metrica sync task",
			zap.Uint("project_id", payload.ProjectID),
			zap.Int("year", payload.Year),
			zap.Int("month", payload.Month),
		)
	}

	// Call sync service method
	err = w.syncService.SyncMetricaData(ctx, payload.ProjectID, payload.Year, payload.Month)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("Failed to sync Metrica data",
				zap.Uint("project_id", payload.ProjectID),
				zap.Error(err),
			)
		}
		return err
	}

	if logger.Log != nil {
		logger.Log.Info("Metrica sync task completed",
			zap.Uint("project_id", payload.ProjectID),
		)
	}

	return nil
}

// handleSyncDirect handles Direct sync task
func (w *Worker) handleSyncDirect(ctx context.Context, task *asynq.Task) error {
	payload, err := ParseSyncDirectPayload(task)
	if err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	if logger.Log != nil {
		logger.Log.Info("Processing Direct sync task",
			zap.Uint("project_id", payload.ProjectID),
			zap.Int("year", payload.Year),
			zap.Int("month", payload.Month),
		)
	}

	// Call sync service method
	err = w.syncService.SyncDirectData(ctx, payload.ProjectID, payload.Year, payload.Month)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("Failed to sync Direct data",
				zap.Uint("project_id", payload.ProjectID),
				zap.Error(err),
			)
		}
		return err
	}

	if logger.Log != nil {
		logger.Log.Info("Direct sync task completed",
			zap.Uint("project_id", payload.ProjectID),
		)
	}

	return nil
}

// handleSyncProject handles project sync task
func (w *Worker) handleSyncProject(ctx context.Context, task *asynq.Task) error {
	payload, err := ParseSyncProjectPayload(task)
	if err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	if logger.Log != nil {
		logger.Log.Info("Processing project sync task",
			zap.Uint("project_id", payload.ProjectID),
		)
	}

	// Call sync service method
	err = w.syncService.SyncProject(ctx, payload.ProjectID)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("Failed to sync project",
				zap.Uint("project_id", payload.ProjectID),
				zap.Error(err),
			)
		}
		return err
	}

	if logger.Log != nil {
		logger.Log.Info("Project sync task completed",
			zap.Uint("project_id", payload.ProjectID),
		)
	}

	return nil
}

// handleAnalyzeMetrics handles metrics analysis task
func (w *Worker) handleAnalyzeMetrics(ctx context.Context, task *asynq.Task) error {
	payload, err := ParseAnalyzeMetricsPayload(task)
	if err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	if logger.Log != nil {
		logger.Log.Info("Processing metrics analysis task",
			zap.Uint("project_id", payload.ProjectID),
			zap.Strings("periods", payload.Periods),
		)
	}

	// Get channel metrics data
	metricsData, err := w.reportService.GetChannelMetrics(ctx, payload.ProjectID, payload.Periods)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("Failed to get channel metrics",
				zap.Uint("project_id", payload.ProjectID),
				zap.Error(err),
			)
		}
		return fmt.Errorf("failed to get channel metrics: %w", err)
	}

	// Analyze metrics using AI
	_, err = w.reportService.AnalyzeChannelMetrics(ctx, metricsData)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("Failed to analyze metrics",
				zap.Uint("project_id", payload.ProjectID),
				zap.Error(err),
			)
		}
		return fmt.Errorf("failed to analyze metrics: %w", err)
	}

	if logger.Log != nil {
		logger.Log.Info("Metrics analysis task completed",
			zap.Uint("project_id", payload.ProjectID),
		)
	}

	return nil
}

// handleGenerateReport handles report generation task
func (w *Worker) handleGenerateReport(ctx context.Context, task *asynq.Task) error {
	payload, err := ParseGenerateReportPayload(task)
	if err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	if logger.Log != nil {
		logger.Log.Info("Processing report generation task",
			zap.Uint("project_id", payload.ProjectID),
		)
	}

	// Generate report
	report, err := w.reportService.GetReport(ctx, payload.ProjectID)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("Failed to generate report",
				zap.Uint("project_id", payload.ProjectID),
				zap.Error(err),
			)
		}
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Store report in cache for retrieval (using project_id as key)
	// TTL: 1 hour - reports are regenerated periodically
	if w.cache != nil {
		cacheKey := fmt.Sprintf("report:project:%d", payload.ProjectID)
		if err := w.cache.Set(cacheKey, report, time.Hour); err != nil {
			if logger.Log != nil {
				logger.Log.Warn("Failed to cache report",
					zap.Uint("project_id", payload.ProjectID),
					zap.Error(err),
				)
			}
		}
	}

	if logger.Log != nil {
		logger.Log.Info("Report generation task completed",
			zap.Uint("project_id", payload.ProjectID),
		)
	}

	return nil
}

// Start starts the worker server
func (w *Worker) Start() error {
	return w.server.Start(w.mux)
}

// Shutdown gracefully shuts down the worker
func (w *Worker) Shutdown() {
	w.server.Shutdown()
}
