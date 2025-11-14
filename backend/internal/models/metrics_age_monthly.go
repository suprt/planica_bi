package models

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
	ID                    uint     `gorm:"primaryKey"`
	ProjectID             uint     `gorm:"not null;index"`
	Year                  int      `gorm:"not null;index"`
	Month                 int      `gorm:"not null;index"`
	AgeGroup              AgeGroup `gorm:"type:enum('18-24','25-34','35-44','45-54','55+','unknown')"`
	Visits                int
	Users                 int
	BounceRate            float64 `gorm:"type:decimal(5,2)"`
	AvgSessionDurationSec int
	CreatedAt             int64
}
