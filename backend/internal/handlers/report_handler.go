package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/cache"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/queue"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/services"
	"go.uber.org/zap"
)

// ReportServiceInterface defines methods for report operations
type ReportServiceInterface interface {
	GetReport(ctx context.Context, projectID uint) (*services.Report, error)
	GetChannelMetrics(ctx context.Context, projectID uint, periods []string) (*services.ChannelMetricsOutput, error)
	AnalyzeChannelMetrics(ctx context.Context, metricsData *services.ChannelMetricsOutput) (*services.MetricsAnalysisResult, error)
	CalculateDynamics(current, previous float64) float64
}

// ReportHandler handles HTTP requests for reports
type ReportHandler struct {
	reportService  ReportServiceInterface
	projectService ProjectServiceInterface
	queueClient    *queue.Client
	cache          *cache.Cache
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService ReportServiceInterface, queueClient *queue.Client, cacheClient *cache.Cache) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		queueClient:   queueClient,
		cache:         cacheClient,
	}
}

// SetProjectService sets the project service (for public reports)
func (h *ReportHandler) SetProjectService(projectService ProjectServiceInterface) {
	h.projectService = projectService
}

// GetReport handles GET /api/report/:id
// Returns JSON with report data for 3 months (M, M-1, M-2)
// If report is not in cache, enqueues generation task and returns task_id
func (h *ReportHandler) GetReport(c echo.Context) error {
	// Force log immediately - this should always appear
	if logger.Log != nil {
		logger.Log.Info("[REPORT HANDLER] GetReport ENTRY",
			zap.String("path", c.Path()),
			zap.String("param_id", c.Param("id")),
		)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("[REPORT HANDLER] ERROR parsing id", zap.Error(err))
		}
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	projectID := uint(id)
	if logger.Log != nil {
		logger.Log.Info("[REPORT HANDLER] GetReport processing",
			zap.Uint("project_id", projectID),
			zap.Bool("cache_is_nil", h.cache == nil),
			zap.Bool("queue_client_is_nil", h.queueClient == nil),
		)
	}

	// Try to get report from cache first
	cacheKey := fmt.Sprintf("report:project:%d", projectID)
	if h.cache != nil {
		var report services.Report
		err := h.cache.Get(cacheKey, &report)
		if err == nil {
			// Report found in cache, return it
			logger.Log.Info("Report retrieved from cache",
				zap.Uint("project_id", projectID),
				zap.String("cache_key", cacheKey),
			)
			return c.JSON(200, report)
		}
		// Cache miss - will enqueue task
		logger.Log.Info("Report not in cache, enqueuing generation task",
			zap.Uint("project_id", projectID),
			zap.String("cache_key", cacheKey),
			zap.String("error", err.Error()),
		)
	} else {
		logger.Log.Warn("Cache is nil, enqueuing generation task",
			zap.Uint("project_id", projectID),
		)
	}

	// Report not in cache, enqueue generation task
	taskInfo, err := h.queueClient.EnqueueGenerateReportTask(projectID)
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Failed to enqueue report generation task: %v", err))
	}

	if taskInfo == nil {
		return echo.NewHTTPError(500, "Failed to enqueue report generation task: no task info returned")
	}

	// Return task info - client should poll or wait for completion
	return c.JSON(202, map[string]interface{}{
		"message":    "Report generation task enqueued",
		"project_id": projectID,
		"task_id":    taskInfo.ID,
		"queue":      taskInfo.Queue,
		"status":     "pending",
		"note":       "Report is being generated. Poll this endpoint or check task status.",
	})
}

// GetPublicReport handles GET /api/public/report/:token
// Returns JSON with report data for 3 months without authentication
func (h *ReportHandler) GetPublicReport(c echo.Context) error {
	ctx := c.Request().Context()

	token := c.Param("token")
	if token == "" {
		return echo.NewHTTPError(400, "Token is required")
	}

	// Get project by public token
	if h.projectService == nil {
		return echo.NewHTTPError(500, "Project service not configured")
	}

	project, err := h.projectService.GetProjectByPublicToken(ctx, token)
	if err != nil {
		return echo.NewHTTPError(404, "Project not found or token invalid")
	}
	if project == nil {
		return echo.NewHTTPError(404, "Project not found")
	}

	// Get report for this project
	report, err := h.reportService.GetReport(ctx, project.ID)
	if err != nil {
		return err
	}

	return c.JSON(200, report)
}

// GetChannelMetrics handles GET /api/channel-metrics/:id?periods=2025-08,2024-09,2024-10
// Returns JSON with channel metrics data from database
func (h *ReportHandler) GetChannelMetrics(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	// Get periods from query parameter
	periodsStr := c.QueryParam("periods")
	if periodsStr == "" {
		return echo.NewHTTPError(400, "periods query parameter is required (comma-separated, e.g., 2025-08,2024-09,2024-10)")
	}

	// Parse periods
	periods := strings.Split(periodsStr, ",")
	for i, period := range periods {
		periods[i] = strings.TrimSpace(period)
	}

	output, err := h.reportService.GetChannelMetrics(ctx, uint(id), periods)
	if err != nil {
		return err
	}

	return c.JSON(200, output)
}

// AnalyzeChannelMetrics handles GET /api/channel-metrics/:id/analyze
// Enqueues a task to analyze channel metrics using AI (asynchronous)
func (h *ReportHandler) AnalyzeChannelMetrics(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	// Get periods from query parameter
	periodsStr := c.QueryParam("periods")
	if periodsStr == "" {
		return echo.NewHTTPError(400, "periods query parameter is required (comma-separated, e.g., 2025-08,2024-09,2024-10)")
	}

	// Parse periods
	periods := strings.Split(periodsStr, ",")
	for i, period := range periods {
		periods[i] = strings.TrimSpace(period)
	}

	// Enqueue analysis task
	taskInfo, err := h.queueClient.EnqueueAnalyzeMetricsTask(uint(id), periods)
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Failed to enqueue analysis task: %v", err))
	}

	if taskInfo == nil {
		return echo.NewHTTPError(500, "Failed to enqueue analysis task: no task info returned")
	}

	return c.JSON(202, map[string]interface{}{
		"message":    "Analysis task enqueued",
		"project_id": uint(id),
		"task_id":    taskInfo.ID,
		"queue":      taskInfo.Queue,
		"status":     "pending",
	})
}
