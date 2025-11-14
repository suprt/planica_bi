package handlers

import (
	"net/http"
)

// DirectHandler handles HTTP requests for Direct accounts
type DirectHandler struct {
	// TODO: add service
}

// NewDirectHandler creates a new Direct handler
func NewDirectHandler() *DirectHandler {
	return &DirectHandler{}
}

// AddDirectAccount handles POST /api/projects/{id}/direct-accounts
func (h *DirectHandler) AddDirectAccount(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

// GetDirectAccounts handles GET /api/projects/{id}/direct-accounts
func (h *DirectHandler) GetDirectAccounts(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

