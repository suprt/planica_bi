package models

import "time"

// Project represents a client project
type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:text;charset=utf8mb4;collate=utf8mb4_unicode_ci;not null" json:"name"`
	Slug        string    `gorm:"type:varchar(191);charset=utf8mb4;collate=utf8mb4_unicode_ci;unique;not null" json:"slug"`
	PublicToken string    `gorm:"type:varchar(64);charset=utf8mb4;collate=utf8mb4_unicode_ci;unique;index" json:"public_token"`
	Timezone    string    `gorm:"type:varchar(191);charset=utf8mb4;collate=utf8mb4_unicode_ci;default:Europe/Moscow" json:"timezone"`
	Currency    string    `gorm:"type:enum('RUB');charset=utf8mb4;collate=utf8mb4_unicode_ci;default:'RUB'" json:"currency"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
