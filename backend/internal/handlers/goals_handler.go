package handlers

import (
	"github.com/labstack/echo/v4"
)

// GoalsHandler handles HTTP requests for goals
type GoalsHandler struct {
	// TODO: add service
}

// NewGoalsHandler creates a new goals handler
func NewGoalsHandler() *GoalsHandler {
	return &GoalsHandler{}
}

// AddGoal handles POST /api/projects/:id/goals
func (h *GoalsHandler) AddGoal(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}

// GetGoals handles GET /api/projects/:id/goals
func (h *GoalsHandler) GetGoals(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}

