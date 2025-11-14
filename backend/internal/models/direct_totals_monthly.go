package models

import "time"

// DirectTotalsMonthly represents monthly totals for all Direct campaigns of a project
type DirectTotalsMonthly struct {
	ID          uint      `gorm:"primaryKey"`
	ProjectID   uint      `gorm:"not null;index"`
	Year        int       `gorm:"not null;index"`
	Month       int       `gorm:"not null;index"`
	Impressions int       `gorm:"not null;default:0"`
	Clicks      int       `gorm:"not null;default:0"`
	CTRPct      float64   `gorm:"type:decimal(6,2)"`
	CPC         float64   `gorm:"type:decimal(12,2)"`
	Conversions *int
	CPA         *float64  `gorm:"type:decimal(12,2)"`
	Cost        float64   `gorm:"type:decimal(14,2);not null;default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
