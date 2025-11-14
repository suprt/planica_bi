package handlers

import (
	"net/http"
)

// GoalsHandler handles HTTP requests for goals
type GoalsHandler struct {
	// TODO: add service
}

// NewGoalsHandler creates a new goals handler
func NewGoalsHandler() *GoalsHandler {
	return &GoalsHandler{}
}

// AddGoal handles POST /api/projects/{id}/goals
func (h *GoalsHandler) AddGoal(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

// GetGoals handles GET /api/projects/{id}/goals
func (h *GoalsHandler) GetGoals(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}

