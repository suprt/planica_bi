package handlers

import (
	"net/http"
)

// SyncHandler handles HTTP requests for synchronization
type SyncHandler struct {
	// TODO: add service
}

// NewSyncHandler creates a new sync handler
func NewSyncHandler() *SyncHandler {
	return &SyncHandler{}
}

// SyncProject handles POST /api/sync/{id}
func (h *SyncHandler) SyncProject(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	// Force synchronization for a specific project
	w.WriteHeader(http.StatusNotImplemented)
}

