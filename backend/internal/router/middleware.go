package router

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/suprt/planica_bi/backend/internal/logger"
)

// timeoutMiddleware sets a timeout for request context
// The context will be cancelled if:
// 1. The parent context is cancelled (client disconnected) - immediate cancellation
// 2. The timeout duration is exceeded (30 seconds)
// This ensures that if a client disconnects, we don't continue processing the request
func timeoutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create context with timeout from parent context
			// If parent is already cancelled (client disconnected), this context will be cancelled immediately
			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			// Replace request context with timeout context
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// zapLoggerMiddleware returns a middleware that logs HTTP requests using zap
func zapLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Log request
			duration := time.Since(start)

			fields := []zap.Field{
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("remote_ip", c.RealIP()),
				zap.Int("status", c.Response().Status),
				zap.Duration("latency", duration),
				zap.String("user_agent", c.Request().UserAgent()),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				if logger.Log != nil {
					logger.Log.Error("HTTP request failed", fields...)
				}
			} else {
				if logger.Log != nil {
					logger.Log.Info("HTTP request", fields...)
				}
			}

			return err
		}
	}
}
