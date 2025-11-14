package handlers

import (
	"net/http"
)

// ReportHandler handles HTTP requests for reports
type ReportHandler struct {
	// TODO: add service
}

// NewReportHandler creates a new report handler
func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

// GetReport handles GET /api/report/{id}
func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	// Returns JSON with report data for 3 months (M, M-1, M-2)
	w.WriteHeader(http.StatusNotImplemented)
}

