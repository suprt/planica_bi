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
	reportService  ReportServiceInterface
	projectService ProjectServiceInterface
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService ReportServiceInterface, projectService ProjectServiceInterface) *ReportHandler {
	return &ReportHandler{
		reportService:  reportService,
		projectService: projectService,
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

// GetPublicReport handles GET /api/public/report/{token}
// Returns JSON with report data for 3 months (M, M-1, M-2) without authentication
func (h *ReportHandler) GetPublicReport(c echo.Context) error {
	ctx := c.Request().Context()
	token := c.Param("token")

	if token == "" {
		return echo.NewHTTPError(400, "Token is required")
	}

	// Get project by public token
	project, err := h.projectService.GetProjectByPublicToken(ctx, token)
	if err != nil {
		return echo.NewHTTPError(404, "Report not found")
	}

	// Check if project is active
	if !project.IsActive {
		return echo.NewHTTPError(404, "Report not found")
	}

	// Get report for the project
	report, err := h.reportService.GetReport(ctx, project.ID)
	if err != nil {
		return err
	}

	return c.JSON(200, report)
}
