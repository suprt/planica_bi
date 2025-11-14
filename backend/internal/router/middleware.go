package router

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
)

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
