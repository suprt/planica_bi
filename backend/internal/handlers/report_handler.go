package handlers

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/services"
)

// ReportServiceInterface defines methods for report operations
type ReportServiceInterface interface {
	GetReport(ctx context.Context, projectID uint) (*services.Report, error)
	CalculateDynamics(current, previous float64) float64
}

// ReportHandler handles HTTP requests for reports
type ReportHandler struct {
	reportService ReportServiceInterface
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService ReportServiceInterface) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
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
