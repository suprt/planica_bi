package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/cache"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/cron"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/database"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/integrations"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/queue"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/router"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/services"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log, err := logger.Init(cfg.AppEnv, cfg.LogPath)
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

	// Initialize cache
	cacheClient, err := cache.NewCache(cfg)
	if err != nil {
		log.Fatal("Failed to initialize cache", zap.Error(err))
	}
	defer cacheClient.Close()
	log.Info("Cache connection established")

	// Initialize repositories
	projectRepo := repositories.NewProjectRepository(db)
	userRepo := repositories.NewUserRepository(db)
	metricsRepo := repositories.NewMetricsRepository(db)
	directRepo := repositories.NewDirectRepository(db, cacheClient)
	counterRepo := repositories.NewCounterRepository(db, cacheClient) // Cached - only changes on manual admin actions
	goalRepo := repositories.NewGoalRepository(db, cacheClient)
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
	userService := services.NewUserService(userRepo)

	// Initialize queue client
	queueClient, err := queue.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to initialize queue client", zap.Error(err))
	}
	defer queueClient.Close()

	// Initialize queue worker
	worker, err := queue.NewWorker(cfg, syncService, reportService, cacheClient)
	if err != nil {
		log.Fatal("Failed to initialize queue worker", zap.Error(err))
	}

	// Start worker in background
	go func() {
		log.Info("Queue worker starting")
		if err := worker.Start(); err != nil {
			log.Fatal("Queue worker failed to start", zap.Error(err))
		}
	}()

	// Initialize cron scheduler
	scheduler := cron.NewScheduler(queueClient, projectService)
	scheduler.StartDailySync()
	scheduler.StartMonthlyFinalization()
	scheduler.Start()
	defer scheduler.Stop()

	log.Info("Cron scheduler initialized and started")

	// Setup routes (pass queueClient instead of syncService)
	e := router.SetupRoutes(
		cfg,
		projectService,
		reportService,
		queueClient,
		goalService,
		directService,
		counterService,
		authService,
		userService,
		userRepo,
		cacheClient,
	)

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		addr := ":8080"
		log.Info("Server starting", zap.String("address", addr))
		if err := e.Start(addr); err != nil {
			log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	// Shutdown worker
	worker.Shutdown()
	log.Info("Server exited")
}
