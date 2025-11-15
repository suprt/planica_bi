package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/database"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/integrations"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/router"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/services"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log, err := logger.Init(cfg.AppEnv)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	log.Info("Application starting",
		zap.String("env", cfg.AppEnv),
		zap.String("version", "1.0.0"),
	)

	// Initialize database connection
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	log.Info("Database connection established")

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	log.Info("Database migrations completed")

	// Initialize repositories
	projectRepo := repositories.NewProjectRepository(db)
	userRepo := repositories.NewUserRepository(db)
	metricsRepo := repositories.NewMetricsRepository(db)
	directRepo := repositories.NewDirectRepository(db)
	counterRepo := repositories.NewCounterRepository(db)
	goalRepo := repositories.NewGoalRepository(db)
	seoRepo := repositories.NewSEORepository(db)

	// Initialize integration clients
	// Note: OAuth token may be empty initially, clients will handle this
	metricaClient := integrations.NewYandexMetricaClient(cfg.YandexOAuthToken)
	directClient := integrations.NewYandexDirectClient(
		cfg.YandexOAuthToken,
		"", // clientLogin will be set per project
		cfg.YandexDirectSandbox,
	)

	// Initialize services
	projectService := services.NewProjectService(projectRepo)
	reportService := services.NewReportService(metricsRepo, directRepo, seoRepo, projectRepo, cfg)
	syncService := services.NewSyncService(
		projectRepo,
		metricsRepo,
		directRepo,
		counterRepo,
		goalRepo,
		metricaClient,
		directClient,
	)
	goalService := services.NewGoalService(goalRepo, counterRepo)
	directService := services.NewDirectService(directRepo)
	counterService := services.NewCounterService(counterRepo)
	authService := services.NewAuthService(
		userRepo,
		cfg.JWTSecret,
		time.Duration(cfg.JWTExpiry)*time.Hour,
	)

	// Setup routes
	e := router.SetupRoutes(
		cfg,
		projectService,
		reportService,
		syncService,
		goalService,
		directService,
		counterService,
		authService,
		userRepo,
	)

	// Start server in a goroutine
	addr := ":8080"
	go func() {
		log.Info("Server starting", zap.String("address", addr))
		if err := e.Start(addr); err != nil {
			log.Error("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Gracefully shutdown server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
