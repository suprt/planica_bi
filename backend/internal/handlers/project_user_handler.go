package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/suprt/planica_bi/backend/internal/logger"
	"github.com/suprt/planica_bi/backend/internal/middleware"
	"github.com/suprt/planica_bi/backend/internal/services"
	"go.uber.org/zap"
)

// ProjectUserHandler handles HTTP requests for project user roles
type ProjectUserHandler struct {
	userService UserServiceInterface
}

// NewProjectUserHandler creates a new project user handler
func NewProjectUserHandler(userService UserServiceInterface) *ProjectUserHandler {
	return &ProjectUserHandler{
		userService: userService,
	}
}

// GetProjectUsers handles GET /api/projects/:id/users
// Returns all users with roles for a project
func (h *ProjectUserHandler) GetProjectUsers(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse project ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid project ID",
		})
	}

	users, err := h.userService.GetProjectUsers(ctx, uint(projectID))
	if err != nil {
		logger.Log.Error("Failed to get project users", zap.Error(err), zap.Uint64("project_id", projectID))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve project users",
		})
	}

	// React-admin expects { data: [...], total: N }
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  users,
		"total": len(users),
	})
}

// AssignUserRole handles POST /api/projects/:id/users
// Assigns a role to user in project
func (h *ProjectUserHandler) AssignUserRole(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse project ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid project ID",
		})
	}

	var req services.AssignRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	req.ProjectID = uint(projectID)

	// Validate request
	if err := middleware.ValidateRequest(c, &req); err != nil {
		return err
	}

	if err := h.userService.AssignRole(ctx, &req); err != nil {
		logger.Log.Error("Failed to assign role",
			zap.Error(err),
			zap.Uint64("project_id", uint64(req.ProjectID)),
			zap.Uint("user_id", req.UserID),
			zap.String("role", req.Role),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	logger.Log.Info("Role assigned",
		zap.Uint64("project_id", uint64(req.ProjectID)),
		zap.Uint("user_id", req.UserID),
		zap.String("role", req.Role),
	)

	// React-admin expects { data: {...} }
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": map[string]interface{}{
			"user_id":    req.UserID,
			"project_id": uint(projectID),
			"role":       req.Role,
		},
	})
}

// UpdateUserRole handles PUT /api/projects/:id/users/:userId
// Updates user's role in project
func (h *ProjectUserHandler) UpdateUserRole(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse project ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid project ID",
		})
	}

	// Parse user ID
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	var req struct {
		Role string `json:"role" validate:"required,oneof=admin manager client"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := middleware.ValidateRequest(c, &req); err != nil {
		return err
	}

	if err := h.userService.UpdateRole(ctx, uint(userID), uint(projectID), req.Role); err != nil {
		logger.Log.Error("Failed to update role",
			zap.Error(err),
			zap.Uint64("project_id", uint64(projectID)),
			zap.Uint64("user_id", userID),
			zap.String("role", req.Role),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	logger.Log.Info("Role updated",
		zap.Uint64("project_id", projectID),
		zap.Uint64("user_id", userID),
		zap.String("role", req.Role),
	)

	// React-admin expects { data: {...} }
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"user_id":    uint(userID),
			"project_id": uint(projectID),
			"role":       req.Role,
		},
	})
}

// RemoveUserRole handles DELETE /api/projects/:id/users/:userId
// Removes user's role from project
func (h *ProjectUserHandler) RemoveUserRole(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse project ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid project ID",
		})
	}

	// Parse user ID
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	if err := h.userService.RemoveRole(ctx, uint(userID), uint(projectID)); err != nil {
		logger.Log.Error("Failed to remove role",
			zap.Error(err),
			zap.Uint64("project_id", projectID),
			zap.Uint64("user_id", userID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to remove role",
		})
	}

	logger.Log.Info("Role removed",
		zap.Uint64("project_id", projectID),
		zap.Uint64("user_id", userID),
	)

	// React-admin expects { data: {...} } with the deleted item
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"user_id":    uint(userID),
			"project_id": uint(projectID),
		},
	})
}
