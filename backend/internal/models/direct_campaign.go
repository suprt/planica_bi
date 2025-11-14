package models

import "time"

// DirectCampaign represents a Yandex.Direct campaign
type DirectCampaign struct {
	ID              uint      `gorm:"primaryKey"`
	DirectAccountID uint      `gorm:"not null;index"`
	CampaignID      int64     `gorm:"not null"`
	Name            string    // Campaign name from Direct API
	Status          string    // Campaign status from Direct API
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}
