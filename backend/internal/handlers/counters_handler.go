package handlers

import (
	"github.com/labstack/echo/v4"
)

// CountersHandler handles HTTP requests for Yandex counters
type CountersHandler struct {
	// TODO: add service
}

// NewCountersHandler creates a new counters handler
func NewCountersHandler() *CountersHandler {
	return &CountersHandler{}
}

// AddCounter handles POST /api/projects/:id/counters
func (h *CountersHandler) AddCounter(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}

// GetCounters handles GET /api/projects/:id/counters
func (h *CountersHandler) GetCounters(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}

