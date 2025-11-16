package models

import "time"

// DirectCampaign represents a Yandex.Direct campaign
type DirectCampaign struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	DirectAccountID uint      `gorm:"not null;index" json:"direct_account_id"`
	CampaignID      int64     `gorm:"not null" json:"campaign_id"`
	Name            string    `gorm:"type:varchar(255)" json:"name"` // Campaign name from Direct API
	Status          string    `gorm:"type:varchar(50)" json:"status"`  // Campaign status from Direct API
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
