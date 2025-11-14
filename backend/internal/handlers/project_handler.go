package handlers

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
)

// ProjectServiceInterface defines methods for project operations
type ProjectServiceInterface interface {
	CreateProject(ctx context.Context, project *models.Project) error
	GetProject(ctx context.Context, id uint) (*models.Project, error)
	GetAllProjects(ctx context.Context) ([]*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, id uint) error
}

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
	projectService ProjectServiceInterface
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService ProjectServiceInterface) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// CreateProject handles POST /api/projects
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	ctx := c.Request().Context()

	var project models.Project
	if err := c.Bind(&project); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	if err := h.projectService.CreateProject(ctx, &project); err != nil {
		// Validation errors should return 400, other errors will be handled by error handler
		if err.Error() == "name is required" || err.Error() == "slug is required" {
			return echo.NewHTTPError(400, err.Error())
		}
		return err
	}

	return c.JSON(201, project)
}

// GetProject handles GET /api/projects/:id
func (h *ProjectHandler) GetProject(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	project, err := h.projectService.GetProject(ctx, uint(id))
	if err != nil {
		return err
	}

	return c.JSON(200, project)
}

// GetAllProjects handles GET /api/projects
func (h *ProjectHandler) GetAllProjects(c echo.Context) error {
	ctx := c.Request().Context()

	projects, err := h.projectService.GetAllProjects(ctx)
	if err != nil {
		return err
	}

	return c.JSON(200, projects)
}

// UpdateProject handles PUT /api/projects/:id
func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	var project models.Project
	if err := c.Bind(&project); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	project.ID = uint(id)

	if err := h.projectService.UpdateProject(ctx, &project); err != nil {
		// Check if it's a validation error or not found error
		if err.Error() == "name is required" || err.Error() == "slug is required" {
			return echo.NewHTTPError(400, err.Error())
		}
		return err
	}

	return c.JSON(200, project)
}

// DeleteProject handles DELETE /api/projects/:id
func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	if err := h.projectService.DeleteProject(ctx, uint(id)); err != nil {
		return err
	}

	return c.NoContent(204)
}
