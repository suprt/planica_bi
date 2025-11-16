package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"go.uber.org/zap"
)

// ProjectServiceInterface defines methods for project operations
type ProjectServiceInterface interface {
	CreateProject(ctx context.Context, project *models.Project) error
	GetProject(ctx context.Context, id uint) (*models.Project, error)
	GetAllProjects(ctx context.Context, userID uint, isAdmin bool) ([]*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, id uint) error
	GetProjectByPublicToken(ctx context.Context, token string) (*models.Project, error)
}

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
	projectService ProjectServiceInterface
	userRepo       repositories.UserRepositoryInterface
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService ProjectServiceInterface, userRepo repositories.UserRepositoryInterface) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		userRepo:       userRepo,
	}
}

// CreateProject handles POST /api/projects
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	ctx := c.Request().Context()

	var project models.Project
	if err := c.Bind(&project); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	// Log received data for debugging
	if logger.Log != nil {
		logger.Log.Info("Creating project",
			zap.String("name", project.Name),
			zap.String("slug", project.Slug),
		)
	}

	if err := h.projectService.CreateProject(ctx, &project); err != nil {
		// Validation errors should return 400, other errors will be handled by error handler
		if err.Error() == "name is required" || err.Error() == "slug is required" {
			return echo.NewHTTPError(400, err.Error())
		}
		return err
	}

	// React-admin expects { data: {...} }
	return c.JSON(201, map[string]interface{}{
		"data": project,
	})
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

	// React-admin expects { data: {...} }
	return c.JSON(200, map[string]interface{}{
		"data": project,
	})
}

// GetPublicLink handles GET /api/projects/:id/public-link
// Returns public URL for accessing report without authentication
func (h *ProjectHandler) GetPublicLink(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	project, err := h.projectService.GetProject(ctx, uint(id))
	if err != nil {
		return err
	}

	if project.PublicToken == "" {
		return echo.NewHTTPError(404, "Public token not found for this project")
	}

	// Get base URL from request or config
	// Use request host/URL for flexible deployment
	protocol := "http"
	if c.Request().TLS != nil {
		protocol = "https"
	}
	host := c.Request().Host
	if host == "" {
		host = "localhost:8080"
	}
	
	publicURL := fmt.Sprintf("%s://%s/api/public/report/%s", protocol, host, project.PublicToken)

	return c.JSON(200, map[string]interface{}{
		"public_url":  publicURL,
		"public_token": project.PublicToken,
		"project_id":   project.ID,
		"project_name": project.Name,
	})
}

// GetAllProjects handles GET /api/projects
func (h *ProjectHandler) GetAllProjects(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context (set by AuthMiddleware)
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return echo.NewHTTPError(401, "User not authenticated")
	}

	// Check if user is admin
	isAdmin, err := h.userRepo.IsAdmin(ctx, userID)
	if err != nil {
		// If error checking admin status, assume not admin
		isAdmin = false
	}

	projects, err := h.projectService.GetAllProjects(ctx, userID, isAdmin)
	if err != nil {
		return err
	}

	// Check if request is from react-admin (has Accept header or specific query param)
	// For now, return both formats - frontend will handle it
	// React-admin expects { data: [...], total: N }
	// Regular API expects [...]
	// We'll return the array format for compatibility, frontend will handle both
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

	// React-admin expects { data: {...} }
	return c.JSON(200, map[string]interface{}{
		"data": project,
	})
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
