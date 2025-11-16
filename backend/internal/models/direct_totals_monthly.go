package models

import "time"

// DirectTotalsMonthly represents monthly totals for all Direct campaigns of a project
type DirectTotalsMonthly struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"not null;index" json:"project_id"`
	Year        int       `gorm:"not null;index" json:"year"`
	Month       int       `gorm:"not null;index" json:"month"`
	Impressions int       `gorm:"not null;default:0" json:"impressions"`
	Clicks      int       `gorm:"not null;default:0" json:"clicks"`
	CTRPct      float64   `gorm:"type:decimal(6,2)" json:"ctr_pct"`
	CPC         float64   `gorm:"type:decimal(12,2)" json:"cpc"`
	Conversions *int      `json:"conversions"`
	CPA         *float64  `gorm:"type:decimal(12,2)" json:"cpa"`
	Cost        float64   `gorm:"type:decimal(14,2);not null;default:0" json:"cost"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
