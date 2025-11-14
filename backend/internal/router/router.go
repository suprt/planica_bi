package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/handlers"
)

// SetupRoutes configures all API routes using Echo
func SetupRoutes() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler()
	countersHandler := handlers.NewCountersHandler()
	directHandler := handlers.NewDirectHandler()
	goalsHandler := handlers.NewGoalsHandler()
	reportHandler := handlers.NewReportHandler()
	syncHandler := handlers.NewSyncHandler()

	// API group
	api := e.Group("/api")

	// Project routes
	api.POST("/projects", projectHandler.CreateProject)
	api.GET("/projects", projectHandler.GetAllProjects)
	api.GET("/projects/:id", projectHandler.GetProject)
	api.PUT("/projects/:id", projectHandler.UpdateProject)
	api.DELETE("/projects/:id", projectHandler.DeleteProject)

	// Counters routes
	api.POST("/projects/:id/counters", countersHandler.AddCounter)
	api.GET("/projects/:id/counters", countersHandler.GetCounters)

	// Direct routes
	api.POST("/projects/:id/direct-accounts", directHandler.AddDirectAccount)
	api.GET("/projects/:id/direct-accounts", directHandler.GetDirectAccounts)

	// Goals routes
	api.POST("/projects/:id/goals", goalsHandler.AddGoal)
	api.GET("/projects/:id/goals", goalsHandler.GetGoals)

	// Report routes
	api.GET("/report/:id", reportHandler.GetReport)

	// Sync routes
	api.POST("/sync/:id", syncHandler.SyncProject)

	return e
}
