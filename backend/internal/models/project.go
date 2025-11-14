package models

import "time"

// Project represents a client project
type Project struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Slug      string    `gorm:"unique;not null"`
	Timezone  string    `gorm:"default:Europe/Moscow"`
	Currency  string    `gorm:"type:enum('RUB');default:'RUB'"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
