package handlers

import (
	"github.com/labstack/echo/v4"
)

// ReportHandler handles HTTP requests for reports
type ReportHandler struct {
	// TODO: add service
}

// NewReportHandler creates a new report handler
func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

// GetReport handles GET /api/report/:id
func (h *ReportHandler) GetReport(c echo.Context) error {
	// TODO: implement
	// Returns JSON with report data for 3 months (M, M-1, M-2)
	return c.NoContent(501) // Not Implemented
}

