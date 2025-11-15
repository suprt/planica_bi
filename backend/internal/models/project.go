package models

import "time"

// Project represents a client project
type Project struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:text;charset=utf8mb4;collate=utf8mb4_unicode_ci;not null"`
	Slug      string    `gorm:"type:varchar(191);charset=utf8mb4;collate=utf8mb4_unicode_ci;unique;not null"`
	Timezone  string    `gorm:"type:varchar(191);charset=utf8mb4;collate=utf8mb4_unicode_ci;default:Europe/Moscow"`
	Currency  string    `gorm:"type:enum('RUB');charset=utf8mb4;collate=utf8mb4_unicode_ci;default:'RUB'"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
