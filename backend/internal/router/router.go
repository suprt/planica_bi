package router

import (
	"net/http"
	
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/handlers"
)

// SetupRoutes configures all API routes
func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	
	// Initialize handlers
	projectHandler := handlers.NewProjectHandler()
	countersHandler := handlers.NewCountersHandler()
	directHandler := handlers.NewDirectHandler()
	goalsHandler := handlers.NewGoalsHandler()
	reportHandler := handlers.NewReportHandler()
	syncHandler := handlers.NewSyncHandler()
	
	// Project routes
	mux.HandleFunc("POST /api/projects", projectHandler.CreateProject)
	mux.HandleFunc("GET /api/projects", projectHandler.GetAllProjects)
	mux.HandleFunc("GET /api/projects/{id}", projectHandler.GetProject)
	mux.HandleFunc("PUT /api/projects/{id}", projectHandler.UpdateProject)
	mux.HandleFunc("DELETE /api/projects/{id}", projectHandler.DeleteProject)
	
	// Counters routes
	mux.HandleFunc("POST /api/projects/{id}/counters", countersHandler.AddCounter)
	mux.HandleFunc("GET /api/projects/{id}/counters", countersHandler.GetCounters)
	
	// Direct routes
	mux.HandleFunc("POST /api/projects/{id}/direct-accounts", directHandler.AddDirectAccount)
	mux.HandleFunc("GET /api/projects/{id}/direct-accounts", directHandler.GetDirectAccounts)
	
	// Goals routes
	mux.HandleFunc("POST /api/projects/{id}/goals", goalsHandler.AddGoal)
	mux.HandleFunc("GET /api/projects/{id}/goals", goalsHandler.GetGoals)
	
	// Report routes
	mux.HandleFunc("GET /api/report/{id}", reportHandler.GetReport)
	
	// Sync routes
	mux.HandleFunc("POST /api/sync/{id}", syncHandler.SyncProject)
	
	return mux
}

