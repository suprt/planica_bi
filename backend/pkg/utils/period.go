package utils

import "time"

// GetReportPeriods returns the three report periods: M (current), M-1, M-2
// Returns periods in format "YYYY-MM"
func GetReportPeriods() (string, string, string) {
	now := time.Now()

	// Current month (M)
	current := now.Format("2006-01")

	// Previous month (M-1)
	prevMonth := now.AddDate(0, -1, 0)
	previous := prevMonth.Format("2006-01")

	// Two months ago (M-2)
	twoMonthsAgo := now.AddDate(0, -2, 0)
	twoMonthsBack := twoMonthsAgo.Format("2006-01")

	return current, previous, twoMonthsBack
}

// CalculateDynamics calculates percentage change
// Formula: (current - previous) / previous * 100
// If previous = 0, returns +0
func CalculateDynamics(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return ((current - previous) / previous) * 100
}
