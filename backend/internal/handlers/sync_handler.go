package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/queue"
)

// SyncHandler handles HTTP requests for synchronization
type SyncHandler struct {
	queueClient *queue.Client
}

// NewSyncHandler creates a new sync handler
func NewSyncHandler(queueClient *queue.Client) *SyncHandler {
	return &SyncHandler{
		queueClient: queueClient,
	}
}

// SyncProject handles POST /api/sync/:id
// Enqueues a task to synchronize a specific project
func (h *SyncHandler) SyncProject(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	// Enqueue sync task
	taskInfo, err := h.queueClient.EnqueueSyncProjectTask(uint(id))
	if err != nil {
		return echo.NewHTTPError(500, "Failed to enqueue sync task: "+err.Error())
	}

	return c.JSON(200, map[string]interface{}{
		"message":    "Synchronization task enqueued",
		"project_id": uint(id),
		"task_id":    taskInfo.ID,
		"queue":      taskInfo.Queue,
	})
}
