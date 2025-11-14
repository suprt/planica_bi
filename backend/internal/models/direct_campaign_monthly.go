package models

// DirectCampaignMonthly represents monthly aggregated metrics for a Direct campaign
type DirectCampaignMonthly struct {
	ID               uint `gorm:"primaryKey"`
	ProjectID        uint `gorm:"not null;index"`
	DirectCampaignID uint `gorm:"not null;index"`
	Year             int  `gorm:"not null;index"`
	Month            int  `gorm:"not null;index"`
	Impressions      int
	Clicks           int
	CTRPct           float64 `gorm:"type:decimal(6,2)"`
	CPC              float64 `gorm:"type:decimal(12,2)"`
	Conversions      *int
	CPA              *float64 `gorm:"type:decimal(12,2)"`
	Cost             float64  `gorm:"type:decimal(14,2)"`
	CreatedAt        int64
}
