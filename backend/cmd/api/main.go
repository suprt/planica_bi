package main

import (
	"log"
	"net/http"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/router"
)

func main() {
	// Load configuration
	_ = config.Load() // TODO: use config for DB connection, etc.

	// TODO: Initialize database connection

	// Setup routes
	mux := router.SetupRoutes()

	// Start server
	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
