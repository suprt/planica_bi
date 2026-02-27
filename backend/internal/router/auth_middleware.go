package router

import (
	"context"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/suprt/planica_bi/backend/internal/logger"
	"github.com/suprt/planica_bi/backend/internal/services"
)

// Context key types to avoid collisions
type contextKey string

const (
	contextKeyUserID    contextKey = "user_id"
	contextKeyProjectID contextKey = "project_id"
	contextKeyUserRole  contextKey = "user_role"
	contextKeyIsAdmin   contextKey = "is_admin"
)

// AuthMiddleware validates JWT token and sets user ID in context
func AuthMiddleware(authService interface {
	ValidateToken(tokenString string) (uint, error)
}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(401, "Authorization header is required")
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(401, "Invalid authorization header format")
			}

			token := parts[1]

			// Validate token
			userID, err := authService.ValidateToken(token)
			if err != nil {
				if logger.Log != nil {
					logger.Log.Warn("Invalid token",
						zap.String("path", c.Request().URL.Path),
						zap.Error(err),
					)
				}
				return echo.NewHTTPError(401, "Invalid or expired token")
			}

			// Set user ID in context
			ctx := context.WithValue(c.Request().Context(), contextKeyUserID, userID)
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set("user_id", userID)

			return next(c)
		}
	}
}

// RequireProjectRole middleware checks if user has required role on project
// For routes without :id parameter, it only checks if user is admin
// For routes with :id parameter, it checks project-specific access
func RequireProjectRole(userRepo services.UserRepositoryInterface, allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user ID from context (set by AuthMiddleware)
			userID, ok := c.Get("user_id").(uint)
			if !ok {
				return echo.NewHTTPError(401, "User not authenticated")
			}

			ctx := c.Request().Context()

			// Check if user is admin (admin has access to all projects)
			isAdmin, err := userRepo.IsAdmin(ctx, userID)
			if err == nil && isAdmin {
				// Set isAdmin in context for handlers
				ctx = context.WithValue(ctx, contextKeyIsAdmin, true)
				c.SetRequest(c.Request().WithContext(ctx))
				c.Set("is_admin", true)
				return next(c)
			}

			// Get project ID from URL parameter
			projectIDStr := c.Param("id")
			if projectIDStr == "" {
				// No project ID in route - only admin can access
				return echo.NewHTTPError(403, "Access denied: admin role required")
			}

			// Parse project ID
			var projectID uint
			if _, err := fmt.Sscanf(projectIDStr, "%d", &projectID); err != nil {
				return echo.NewHTTPError(400, "Invalid project ID")
			}

			// Get user's role for this project
			role, err := userRepo.GetUserProjectRole(ctx, userID, projectID)
			if err != nil {
				return echo.NewHTTPError(403, "Access denied: no permission for this project")
			}

			// Check if user's role is in allowed roles
			hasAccess := false
			for _, allowedRole := range allowedRoles {
				if role.Role == allowedRole {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				return echo.NewHTTPError(403, "Access denied: insufficient permissions")
			}

			// Set project ID and role in context
			ctx = context.WithValue(ctx, contextKeyProjectID, projectID)
			ctx = context.WithValue(ctx, contextKeyUserRole, role.Role)
			ctx = context.WithValue(ctx, contextKeyIsAdmin, false)
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set("project_id", projectID)
			c.Set("user_role", role.Role)
			c.Set("is_admin", false)

			return next(c)
		}
	}
}
