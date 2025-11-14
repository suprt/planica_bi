package services

// ReportService handles business logic for reports
type ReportService struct {
	// TODO: add repositories
}

// NewReportService creates a new report service
func NewReportService() *ReportService {
	return &ReportService{}
}

// GetReport generates a report for a project for the last 3 months
func (s *ReportService) GetReport(projectID uint) (interface{}, error) {
	// TODO: implement
	// Returns report with periods M, M-1, M-2
	// Includes: metrica summary, age breakdown, direct totals, campaigns, seo
	return nil, nil
}

// CalculateDynamics calculates percentage change between two values
func (s *ReportService) CalculateDynamics(current, previous float64) float64 {
	// TODO: implement
	// Formula: (current - previous) / previous * 100
	// If previous = 0, return +0
	return 0
}

