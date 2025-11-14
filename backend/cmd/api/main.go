package main

import (
	"log"
	"net/http"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/database"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup routes (pass db connection if needed)
	_ = db // TODO: pass db to router/handlers when implementing
	mux := router.SetupRoutes()

	// Start server
	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
