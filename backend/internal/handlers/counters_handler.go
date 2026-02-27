package handlers

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/suprt/planica_bi/backend/internal/models"
)

// CounterServiceInterface defines methods for counter operations
type CounterServiceInterface interface {
	CreateCounter(ctx context.Context, counter *models.YandexCounter) error
	GetCountersByProject(ctx context.Context, projectID uint) ([]*models.YandexCounter, error)
}

// CountersHandler handles HTTP requests for Yandex counters
type CountersHandler struct {
	counterService CounterServiceInterface
}

// NewCountersHandler creates a new counters handler
func NewCountersHandler(counterService CounterServiceInterface) *CountersHandler {
	return &CountersHandler{
		counterService: counterService,
	}
}

// AddCounter handles POST /api/projects/:id/counters
func (h *CountersHandler) AddCounter(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	var counter models.YandexCounter
	if err := c.Bind(&counter); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	// Set project ID from URL
	counter.ProjectID = uint(projectID)

	if err := h.counterService.CreateCounter(ctx, &counter); err != nil {
		// Validation errors should return 400, other errors will be handled by error handler
		if err.Error() == "project_id is required" || err.Error() == "counter_id is required" {
			return echo.NewHTTPError(400, err.Error())
		}
		if err.Error() == "counter with this CounterID already exists" {
			return echo.NewHTTPError(409, err.Error())
		}
		return err
	}

	return c.JSON(201, counter)
}

// GetCounters handles GET /api/projects/:id/counters
func (h *CountersHandler) GetCounters(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	counters, err := h.counterService.GetCountersByProject(ctx, uint(projectID))
	if err != nil {
		return err
	}

	return c.JSON(200, counters)
}
