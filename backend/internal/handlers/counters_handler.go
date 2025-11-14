package handlers

import (
	"net/http"
)

// CountersHandler handles HTTP requests for Yandex counters
type CountersHandler struct {
	// TODO: add service
}

// NewCountersHandler creates a new counters handler
func NewCountersHandler() *CountersHandler {
	return &CountersHandler{}
}

// AddCounter handles POST /api/projects/{id}/counters
func (h *CountersHandler) AddCounter(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

// GetCounters handles GET /api/projects/{id}/counters
func (h *CountersHandler) GetCounters(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

