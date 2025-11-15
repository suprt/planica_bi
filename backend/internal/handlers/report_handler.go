package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/queue"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/services"
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
	reportService ReportServiceInterface
	queueClient   *queue.Client
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService ReportServiceInterface, queueClient *queue.Client) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		queueClient:   queueClient,
	}
}

// GetReport handles GET /api/report/:id
// Returns JSON with report data for 3 months (M, M-1, M-2)
func (h *ReportHandler) GetReport(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	report, err := h.reportService.GetReport(ctx, uint(id))
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

	return c.JSON(202, map[string]interface{}{
		"message":    "Analysis task enqueued",
		"project_id": uint(id),
		"task_id":    taskInfo.ID,
		"queue":      taskInfo.Queue,
		"status":    "pending",
	})
}
