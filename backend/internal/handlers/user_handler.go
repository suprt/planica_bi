package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/services"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for user management
type UserHandler struct {
	userService UserServiceInterface
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAllUsers handles GET /api/users
// Returns list of all users (admin only)
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		logger.Log.Error("Failed to get users", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve users",
		})
	}

	// React-admin expects { data: [...], total: N }
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  users,
		"total": len(users),
	})
}

// GetUser handles GET /api/users/:id
// Returns single user with their projects
func (h *UserHandler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse user ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	user, err := h.userService.GetUserByID(ctx, uint(userID))
	if err != nil {
		logger.Log.Error("Failed to get user", zap.Error(err), zap.Uint64("user_id", userID))
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// React-admin expects { data: {...} }
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": user,
	})
}

// CreateUser handles POST /api/users
// Creates a new user (admin only)
func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req services.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		logger.Log.Error("Failed to create user", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	logger.Log.Info("User created", zap.Uint("user_id", user.ID), zap.String("email", user.Email))

	// React-admin expects { data: {...} }
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": user,
	})
}

// UpdateUser handles PUT /api/users/:id
// Updates user information (admin only)
func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse user ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	var req services.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.UpdateUser(ctx, uint(userID), &req)
	if err != nil {
		logger.Log.Error("Failed to update user", zap.Error(err), zap.Uint64("user_id", userID))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	logger.Log.Info("User updated", zap.Uint("user_id", user.ID))

	// React-admin expects { data: {...} }
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": user,
	})
}

// DeleteUser handles DELETE /api/users/:id
// Deletes a user (admin only)
func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse user ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	if err := h.userService.DeleteUser(ctx, uint(userID)); err != nil {
		logger.Log.Error("Failed to delete user", zap.Error(err), zap.Uint64("user_id", userID))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete user",
		})
	}

	logger.Log.Info("User deleted", zap.Uint64("user_id", userID))

	// React-admin expects { data: {...} } with the deleted item
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]uint64{"id": userID},
	})
}

// GetUserProjects handles GET /api/users/:id/projects
// Returns all projects for a user
func (h *UserHandler) GetUserProjects(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse user ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	user, err := h.userService.GetUserByID(ctx, uint(userID))
	if err != nil {
		logger.Log.Error("Failed to get user projects", zap.Error(err), zap.Uint64("user_id", userID))
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  user.Projects,
		"total": len(user.Projects),
	})
}

