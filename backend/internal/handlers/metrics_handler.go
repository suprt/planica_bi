package handlers

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/suprt/planica_bi/backend/internal/services"
)

// MetricsServiceInterface defines methods for Metrics operations
type MetricsServiceInterface interface {
	GetMetricsWithData(ctx context.Context, projectID uint) (*services.MetricsWithData, error)
}

// MetricsHandler handles HTTP requests for Metrics
type MetricsHandler struct {
	metricsService MetricsServiceInterface
}

// NewMetricsHandler creates a new Metrics handler
func NewMetricsHandler(metricsService MetricsServiceInterface) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
	}
}

// GetMetrics handles GET /api/projects/:id/metrics
func (h *MetricsHandler) GetMetrics(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	metrics, err := h.metricsService.GetMetricsWithData(ctx, uint(projectID))
	if err != nil {
		return err
	}

	return c.JSON(200, metrics)
}
