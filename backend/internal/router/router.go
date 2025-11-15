package router

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/handlers"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
)

// SetupRoutes configures all API routes using Echo
// Services are passed as parameters to handlers
func SetupRoutes(
	cfg *config.Config,
	projectService handlers.ProjectServiceInterface,
	reportService handlers.ReportServiceInterface,
	syncService handlers.SyncServiceInterface,
	goalService handlers.GoalServiceInterface,
	directService handlers.DirectServiceInterface,
	counterService handlers.CounterServiceInterface,
	authService handlers.AuthServiceInterface,
	userRepo repositories.UserRepositoryInterface,
) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Set request timeout (30 seconds)
	e.Use(timeoutMiddleware(30 * time.Second))

	// Custom logger middleware using zap
	e.Use(zapLoggerMiddleware())

	// Custom error handler
	e.HTTPErrorHandler = customErrorHandler

	// Initialize handlers with services
	projectHandler := handlers.NewProjectHandler(projectService, userRepo)
	countersHandler := handlers.NewCountersHandler(counterService)
	directHandler := handlers.NewDirectHandler(directService)
	goalsHandler := handlers.NewGoalsHandler(goalService)
	reportHandler := handlers.NewReportHandler(reportService)
	syncHandler := handlers.NewSyncHandler(syncService)
	oauthHandler := handlers.NewOAuthHandler(cfg)
	authHandler := handlers.NewAuthHandler(authService)

	// API group
	api := e.Group("/api")

	// Public routes (no authentication required)
	// Auth routes
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// OAuth routes (public, but may require auth later)
	api.GET("/oauth/yandex", oauthHandler.InitiateAuth)
	api.GET("/oauth/yandex/callback", oauthHandler.HandleCallback)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(AuthMiddleware(authService))

	// Project routes
	// Create project (admin only)
	adminOnly := protected.Group("")
	adminOnly.Use(RequireProjectRole(userRepo, "admin"))
	adminOnly.POST("/projects", projectHandler.CreateProject)
	adminOnly.DELETE("/projects/:id", projectHandler.DeleteProject)
	adminOnly.POST("/sync/:id", syncHandler.SyncProject)

	// Get all projects (users see only their projects - handled in service)
	protected.GET("/projects", projectHandler.GetAllProjects)

	// Project-specific routes (require project access)
	projectRoutes := protected.Group("")
	projectRoutes.Use(RequireProjectRole(userRepo, "admin", "manager", "client"))
	projectRoutes.GET("/projects/:id", projectHandler.GetProject)
	projectRoutes.GET("/projects/:id/counters", countersHandler.GetCounters)
	projectRoutes.GET("/projects/:id/direct-accounts", directHandler.GetDirectAccounts)
	projectRoutes.GET("/projects/:id/goals", goalsHandler.GetGoals)
	projectRoutes.GET("/report/:id", reportHandler.GetReport)
	projectRoutes.GET("/channel-metrics/:id", reportHandler.GetChannelMetrics)
	projectRoutes.GET("/channel-metrics/:id/analyze", reportHandler.AnalyzeChannelMetrics)

	// Manager and admin routes (require manager or admin role)
	managerRoutes := protected.Group("")
	managerRoutes.Use(RequireProjectRole(userRepo, "admin", "manager"))
	managerRoutes.PUT("/projects/:id", projectHandler.UpdateProject)
	managerRoutes.POST("/projects/:id/counters", countersHandler.AddCounter)
	managerRoutes.POST("/projects/:id/direct-accounts", directHandler.AddDirectAccount)
	managerRoutes.POST("/projects/:id/goals", goalsHandler.AddGoal)

	return e
}
