package models

// DirectCampaign represents a Yandex.Direct campaign
type DirectCampaign struct {
	ID              uint  `gorm:"primaryKey"`
	DirectAccountID uint  `gorm:"not null;index"`
	CampaignID      int64 `gorm:"not null"`
	Name            string
	Status          string
	CreatedAt       int64
	UpdatedAt       int64
}
