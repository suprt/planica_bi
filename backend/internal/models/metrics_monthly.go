package models

import "time"

// MetricsMonthly represents monthly aggregated metrics from Yandex.Metrica
type MetricsMonthly struct {
	ID                    uint    `gorm:"primaryKey" json:"id"`
	ProjectID             uint    `gorm:"not null;index" json:"project_id"`
	Year                  int     `gorm:"not null;index" json:"year"`
	Month                 int     `gorm:"not null;index" json:"month"`
	Visits                int     `gorm:"not null;default:0" json:"visits"`
	Users                 int     `gorm:"not null;default:0" json:"users"`
	BounceRate            float64 `gorm:"type:decimal(5,2)" json:"bounce_rate"`
	AvgSessionDurationSec int     `gorm:"not null;default:0" json:"avg_session_duration_sec"`
	Conversions           *int    `json:"conversions"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for MetricsMonthly
func (MetricsMonthly) TableName() string {
	return "metrics_monthly"
}
