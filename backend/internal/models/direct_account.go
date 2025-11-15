package models

import "time"

// DirectAccount represents a Yandex.Direct account linked to a project
type DirectAccount struct {
	ID          uint      `gorm:"primaryKey"`
	ProjectID   uint      `gorm:"not null;index"`
	ClientLogin string    `gorm:"type:varchar(255);not null"`
	AccountName *string   `gorm:"type:varchar(255)"` // Optional account name
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
