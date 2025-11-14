package utils

import "time"

// GetMonthStart returns the first day of the current month
func GetMonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// GetMonthEnd returns the last day of the current month
func GetMonthEnd(t time.Time) time.Time {
	nextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return nextMonth.AddDate(0, 0, -1)
}

// FormatPeriod formats year and month as "YYYY-MM"
func FormatPeriod(year, month int) string {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01")
}

