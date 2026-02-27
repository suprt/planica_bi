package router

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/suprt/planica_bi/backend/internal/cache"
	"github.com/suprt/planica_bi/backend/internal/config"
	"github.com/suprt/planica_bi/backend/internal/handlers"
	"github.com/suprt/planica_bi/backend/internal/middleware"
	"github.com/suprt/planica_bi/backend/internal/queue"
	"github.com/suprt/planica_bi/backend/internal/services"
)

// Router holds the Echo instance and middleware
type Router struct {
	Echo        *echo.Echo
	rateLimiter *middleware.RateLimiter
	authLimiter *middleware.RateLimiter
}

// SetupRoutes configures all API routes using Echo
// Services are passed as parameters to handlers
func SetupRoutes(
	cfg *config.Config,
	projectService handlers.ProjectServiceInterface,
	reportService handlers.ReportServiceInterface,
	queueClient *queue.Client,
	goalService handlers.GoalServiceInterface,
	directService handlers.DirectServiceInterface,
	counterService handlers.CounterServiceInterface,
	metricsService handlers.MetricsServiceInterface,
	marketingService handlers.MarketingServiceInterface,
	authService handlers.AuthServiceInterface,
	userService handlers.UserServiceInterface,
	userRepo services.UserRepositoryInterface,
	cacheClient *cache.Cache,
) *Router {
	router := &Router{
		Echo: echo.New(),
	}
	e := router.Echo

	// Middleware
	e.Use(echoMiddleware.Recover())

	// CORS with configuration
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{cfg.FrontendURL},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Rate limiting - general limiter for all routes
	router.rateLimiter = middleware.NewRateLimiter(middleware.DefaultRateLimiterConfig())
	e.Use(router.rateLimiter.Middleware())

	// UTF-8 encoding middleware for responses
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set UTF-8 charset for JSON responses
			c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
			return next(c)
		}
	})

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
	metricsHandler := handlers.NewMetricsHandler(metricsService)
	reportHandler := handlers.NewReportHandler(reportService, queueClient, cacheClient)
	reportHandler.SetProjectService(projectService) // Set project service for public reports
	syncHandler := handlers.NewSyncHandler(queueClient)
	oauthHandler := handlers.NewOAuthHandler(cfg)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	projectUserHandler := handlers.NewProjectUserHandler(userService)
	healthHandler := handlers.NewHealthHandler()

	// Health check routes (public, no authentication required)
	e.GET("/health", healthHandler.Health)
	e.GET("/ready", healthHandler.Ready)

	// API group
	api := e.Group("/api")

	// Strict rate limiter for auth endpoints (prevent brute force)
	router.authLimiter = middleware.NewRateLimiter(middleware.StrictRateLimiterConfig())

	// Public routes (no authentication required)
	// Auth routes with strict rate limiting
	api.POST("/auth/register", authHandler.Register, router.authLimiter.Middleware())
	api.POST("/auth/login", authHandler.Login, router.authLimiter.Middleware())

	// OAuth routes (public, but may require auth later)
	api.GET("/oauth/yandex", oauthHandler.InitiateAuth)
	api.GET("/oauth/yandex/callback", oauthHandler.HandleCallback)

	// Public report route (no authentication required)
	api.GET("/public/report/:token", reportHandler.GetPublicReport)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(AuthMiddleware(authService))

	// OAuth status (protected - requires authentication)
	protected.GET("/oauth/status", oauthHandler.GetOAuthStatus)

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
	projectRoutes.GET("/projects/:id/public-link", projectHandler.GetPublicLink)
	projectRoutes.GET("/projects/:id/counters", countersHandler.GetCounters)
	projectRoutes.GET("/projects/:id/direct-accounts", directHandler.GetDirectAccounts)
	projectRoutes.GET("/projects/:id/campaigns", directHandler.GetCampaigns)
	projectRoutes.GET("/projects/:id/metrics", metricsHandler.GetMetrics)
	projectRoutes.GET("/projects/:id/marketing", handlers.NewMarketingHandler(marketingService).GetMarketing)
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

	// Admin panel routes (require admin role)
	// User management
	adminOnly.GET("/users", userHandler.GetAllUsers)
	adminOnly.GET("/users/:id", userHandler.GetUser)
	adminOnly.POST("/users", userHandler.CreateUser)
	adminOnly.PUT("/users/:id", userHandler.UpdateUser)
	adminOnly.DELETE("/users/:id", userHandler.DeleteUser)
	adminOnly.GET("/users/:id/projects", userHandler.GetUserProjects)

	// Project user role management
	adminOnly.GET("/projects/:id/users", projectUserHandler.GetProjectUsers)
	adminOnly.POST("/projects/:id/users", projectUserHandler.AssignUserRole)
	adminOnly.PUT("/projects/:id/users/:userId", projectUserHandler.UpdateUserRole)
	adminOnly.DELETE("/projects/:id/users/:userId", projectUserHandler.RemoveUserRole)

	return router
}

// Shutdown stops all background processes and cleanup resources
func (r *Router) Shutdown(ctx context.Context) error {
	if r.rateLimiter != nil {
		r.rateLimiter.Stop()
	}
	if r.authLimiter != nil {
		r.authLimiter.Stop()
	}
	return nil
}
