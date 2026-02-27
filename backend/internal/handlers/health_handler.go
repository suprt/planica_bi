package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthHandler provides health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Health returns the health status of the application
func (h *HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
	})
}

// Ready returns the readiness status of the application
func (h *HealthHandler) Ready(c echo.Context) error {
	// Add more sophisticated checks here (database, cache, etc.)
	return c.JSON(http.StatusOK, HealthResponse{
		Status:  "ready",
		Version: "1.0.0",
	})
}
