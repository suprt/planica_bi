package models

import "time"

// AgeGroup represents age group enum
type AgeGroup string

const (
	AgeGroup1824    AgeGroup = "18-24"
	AgeGroup2534    AgeGroup = "25-34"
	AgeGroup3544    AgeGroup = "35-44"
	AgeGroup4554    AgeGroup = "45-54"
	AgeGroup55Plus  AgeGroup = "55+"
	AgeGroupUnknown AgeGroup = "unknown"
)

// MetricsAgeMonthly represents monthly metrics broken down by age group
type MetricsAgeMonthly struct {
	ID                    uint      `gorm:"primaryKey"`
	ProjectID             uint      `gorm:"not null;index"`
	Year                  int       `gorm:"not null;index"`
	Month                 int       `gorm:"not null;index"`
	AgeGroup              AgeGroup  `gorm:"type:varchar(20);not null"`
	Visits                int       `gorm:"not null;default:0"`
	Users                 int       `gorm:"not null;default:0"`
	BounceRate            float64   `gorm:"type:decimal(5,2)"`
	AvgSessionDurationSec int       `gorm:"not null;default:0"`
	CreatedAt             time.Time `gorm:"autoCreateTime"`
}

// TableName specifies the table name for MetricsAgeMonthly
func (MetricsAgeMonthly) TableName() string {
	return "metrics_age_monthly"
}
