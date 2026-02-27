package router

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/suprt/planica_bi/backend/internal/logger"
)

// customErrorHandler handles all errors returned by handlers
func customErrorHandler(err error, c echo.Context) {
	// Default status code
	code := http.StatusInternalServerError
	message := "Internal server error"

	// Check if it's an HTTP error from Echo
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		} else {
			message = http.StatusText(code)
		}
	} else {
		// Check for common errors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
			message = "Resource not found"
		} else if errors.Is(err, context.DeadlineExceeded) {
			// Request timeout
			code = http.StatusRequestTimeout
			message = "Request timeout"
			if logger.Log != nil {
				logger.Log.Warn("Request timeout",
					zap.String("path", c.Request().URL.Path),
					zap.String("method", c.Request().Method),
				)
			}
		} else if errors.Is(err, context.Canceled) {
			// Client cancelled the request
			code = 499 // Client Closed Request (non-standard but commonly used)
			message = "Client closed request"
			if logger.Log != nil {
				logger.Log.Info("Client cancelled request",
					zap.String("path", c.Request().URL.Path),
					zap.String("method", c.Request().Method),
				)
			}
		} else if isDuplicateEntryError(err) {
			code = http.StatusConflict
			message = "Resource already exists"
		} else {
			// Log unexpected errors
			if logger.Log != nil {
				logger.Log.Error("Unexpected error",
					zap.Error(err),
					zap.String("path", c.Request().URL.Path),
					zap.String("method", c.Request().Method),
				)
			}
			message = err.Error()
		}
	}

	// Log error
	if logger.Log != nil {
		fields := []zap.Field{
			zap.Int("status", code),
			zap.String("path", c.Request().URL.Path),
			zap.String("method", c.Request().Method),
		}

		if code >= 500 {
			logger.Log.Error("HTTP error", append(fields, zap.Error(err))...)
		} else {
			logger.Log.Warn("HTTP error", append(fields, zap.String("message", message))...)
		}
	}

	// Send JSON response
	if !c.Response().Committed {
		if err := c.JSON(code, map[string]interface{}{
			"error": message,
		}); err != nil {
			// If JSON encoding fails, log it
			if logger.Log != nil {
				logger.Log.Error("Failed to send error response", zap.Error(err))
			}
		}
	}
}

// isDuplicateEntryError checks if the error is a MySQL duplicate entry error
func isDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	// Check for MySQL duplicate entry error patterns
	return strings.Contains(errStr, "duplicate entry") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "1062") // MySQL error code for duplicate entry
}
