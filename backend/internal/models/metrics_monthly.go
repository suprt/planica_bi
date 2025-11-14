package models

import "time"

// MetricsMonthly represents monthly aggregated metrics from Yandex.Metrica
type MetricsMonthly struct {
	ID                    uint      `gorm:"primaryKey"`
	ProjectID             uint      `gorm:"not null;index"`
	Year                  int       `gorm:"not null;index"`
	Month                 int       `gorm:"not null;index"`
	Visits                int       `gorm:"not null;default:0"`
	Users                 int       `gorm:"not null;default:0"`
	BounceRate            float64   `gorm:"type:decimal(5,2)"`
	AvgSessionDurationSec int       `gorm:"not null;default:0"`
	Conversions           *int
	CreatedAt             time.Time `gorm:"autoCreateTime"`
}

