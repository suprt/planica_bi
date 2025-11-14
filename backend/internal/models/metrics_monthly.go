package models

// MetricsMonthly represents monthly aggregated metrics from Yandex.Metrica
type MetricsMonthly struct {
	ID                      uint    `gorm:"primaryKey"`
	ProjectID               uint    `gorm:"not null;index"`
	Year                    int     `gorm:"not null;index"`
	Month                   int     `gorm:"not null;index"`
	Visits                  int
	Users                   int
	BounceRate              float64 `gorm:"type:decimal(5,2)"`
	AvgSessionDurationSec   int
	Conversions             *int
	CreatedAt               int64
}

