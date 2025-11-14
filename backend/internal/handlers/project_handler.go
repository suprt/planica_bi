package handlers

import (
	"github.com/labstack/echo/v4"
)

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
	// TODO: add service
}

// NewProjectHandler creates a new project handler
func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{}
}

// CreateProject handles POST /api/projects
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}

// GetProject handles GET /api/projects/:id
func (h *ProjectHandler) GetProject(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}

// GetAllProjects handles GET /api/projects
func (h *ProjectHandler) GetAllProjects(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}

// UpdateProject handles PUT /api/projects/:id
func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}

// DeleteProject handles DELETE /api/projects/:id
func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}
