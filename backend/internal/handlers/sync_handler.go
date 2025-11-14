package handlers

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
)

// SyncServiceInterface defines methods for synchronization operations
type SyncServiceInterface interface {
	SyncProject(ctx context.Context, projectID uint) error
	SyncAllProjects(ctx context.Context) error
	FinalizeMonth(ctx context.Context) error
}

// SyncHandler handles HTTP requests for synchronization
type SyncHandler struct {
	syncService SyncServiceInterface
}

// NewSyncHandler creates a new sync handler
func NewSyncHandler(syncService SyncServiceInterface) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
	}
}

// SyncProject handles POST /api/sync/:id
// Force synchronization for a specific project
func (h *SyncHandler) SyncProject(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	if err := h.syncService.SyncProject(ctx, uint(id)); err != nil {
		return err
	}

	return c.JSON(200, map[string]interface{}{
		"message":    "Synchronization started",
		"project_id": uint(id),
	})
}
