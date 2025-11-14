package main

import (
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/database"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/router"
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

	// TODO: Initialize repositories and services
	// For now, passing nil - will be implemented when services are ready
	// projectRepo := repositories.NewProjectRepository(db)
	// projectService := services.NewProjectService(projectRepo)
	// reportService := services.NewReportService(...)
	// syncService := services.NewSyncService(...)
	// goalService := services.NewGoalService(...)
	// directService := services.NewDirectService(...)
	// counterService := services.NewCounterService(...)

	_ = db                                                // Will be used when initializing repositories
	e := router.SetupRoutes(nil, nil, nil, nil, nil, nil) // TODO: pass actual services

	// Start server
	addr := ":8080"
	log.Info("Server starting", zap.String("address", addr))
	if err := e.Start(addr); err != nil {
		log.Fatal("Server failed to start", zap.Error(err))
	}
}
