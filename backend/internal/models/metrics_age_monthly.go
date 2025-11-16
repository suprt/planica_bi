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
	ID                    uint      `gorm:"primaryKey" json:"id"`
	ProjectID             uint      `gorm:"not null;index" json:"project_id"`
	Year                  int       `gorm:"not null;index" json:"year"`
	Month                 int       `gorm:"not null;index" json:"month"`
	AgeGroup              AgeGroup  `gorm:"type:varchar(20);not null" json:"age_group"`
	Visits                int       `gorm:"not null;default:0" json:"visits"`
	Users                 int       `gorm:"not null;default:0" json:"users"`
	BounceRate            float64   `gorm:"type:decimal(5,2)" json:"bounce_rate"`
	AvgSessionDurationSec int       `gorm:"not null;default:0" json:"avg_session_duration_sec"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for MetricsAgeMonthly
func (MetricsAgeMonthly) TableName() string {
	return "metrics_age_monthly"
}
