package handlers

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/suprt/planica_bi/backend/internal/models"
)

// GoalServiceInterface defines methods for goal operations
type GoalServiceInterface interface {
	CreateGoal(ctx context.Context, goal *models.Goal) error
	GetGoalsByProject(ctx context.Context, projectID uint) ([]*models.Goal, error)
}

// GoalsHandler handles HTTP requests for goals
type GoalsHandler struct {
	goalService GoalServiceInterface
}

// NewGoalsHandler creates a new goals handler
func NewGoalsHandler(goalService GoalServiceInterface) *GoalsHandler {
	return &GoalsHandler{
		goalService: goalService,
	}
}

// AddGoal handles POST /api/projects/:id/goals
// Creates a new goal for a counter (counterID must be provided in request body)
func (h *GoalsHandler) AddGoal(c echo.Context) error {
	ctx := c.Request().Context()

	_, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	var goal models.Goal
	if err := c.Bind(&goal); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	if err := h.goalService.CreateGoal(ctx, &goal); err != nil {
		// Validation errors should return 400, other errors will be handled by error handler
		if err.Error() == "counter_id is required" || err.Error() == "goal_id is required" {
			return echo.NewHTTPError(400, err.Error())
		}
		if err.Error() == "counter not found" {
			return echo.NewHTTPError(400, err.Error())
		}
		if err.Error() == "goal with this GoalID already exists for this counter" {
			return echo.NewHTTPError(409, err.Error())
		}
		return err
	}

	return c.JSON(201, goal)
}

// GetGoals handles GET /api/projects/:id/goals
// Returns all goals for all counters of the project
func (h *GoalsHandler) GetGoals(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	goals, err := h.goalService.GetGoalsByProject(ctx, uint(projectID))
	if err != nil {
		return err
	}

	return c.JSON(200, goals)
}
