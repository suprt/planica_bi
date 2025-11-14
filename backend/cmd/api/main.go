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

	// Setup routes (pass db connection if needed)
	_ = db // TODO: pass db to router/handlers when implementing
	e := router.SetupRoutes()

	// Start server
	addr := ":8080"
	log.Info("Server starting", zap.String("address", addr))
	if err := e.Start(addr); err != nil {
		log.Fatal("Server failed to start", zap.Error(err))
	}
}
