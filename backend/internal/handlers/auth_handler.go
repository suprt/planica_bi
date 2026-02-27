package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/suprt/planica_bi/backend/internal/middleware"
	"github.com/suprt/planica_bi/backend/internal/services"
)

// AuthServiceInterface defines methods for authentication operations
type AuthServiceInterface interface {
	Register(ctx context.Context, req *services.RegisterRequest) (*services.AuthResponse, error)
	Login(ctx context.Context, req *services.LoginRequest) (*services.AuthResponse, error)
	ValidateToken(tokenString string) (uint, error)
}

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService AuthServiceInterface
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()

	var req services.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	// Validate request
	if err := middleware.ValidateRequest(c, &req); err != nil {
		return err
	}

	response, err := h.authService.Register(ctx, &req)
	if err != nil {
		// Check for validation errors
		if err.Error() == "name is required" ||
			err.Error() == "email is required" ||
			err.Error() == "password is required" ||
			err.Error() == "password must be at least 8 characters" ||
			err.Error() == "user with this email already exists" {
			return echo.NewHTTPError(400, err.Error())
		}
		return err
	}

	return c.JSON(201, response)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var req services.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	// Validate request
	if err := middleware.ValidateRequest(c, &req); err != nil {
		return err
	}

	response, err := h.authService.Login(ctx, &req)
	if err != nil {
		// Check for validation errors
		if err.Error() == "email is required" ||
			err.Error() == "password is required" ||
			err.Error() == "invalid email or password" ||
			err.Error() == "user account is inactive" {
			return echo.NewHTTPError(401, err.Error())
		}
		return err
	}

	return c.JSON(200, response)
}
