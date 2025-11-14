package models

import "time"

// DirectAccount represents a Yandex.Direct account linked to a project
type DirectAccount struct {
	ID          uint      `gorm:"primaryKey"`
	ProjectID   uint      `gorm:"not null;index"`
	ClientLogin string    `gorm:"not null"`
	AccountName string    // Optional account name
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
