package handlers

import (
	"github.com/labstack/echo/v4"
)

// SyncHandler handles HTTP requests for synchronization
type SyncHandler struct {
	// TODO: add service
}

// NewSyncHandler creates a new sync handler
func NewSyncHandler() *SyncHandler {
	return &SyncHandler{}
}

// SyncProject handles POST /api/sync/:id
func (h *SyncHandler) SyncProject(c echo.Context) error {
	// TODO: implement
	// Force synchronization for a specific project
	return c.NoContent(501) // Not Implemented
}

