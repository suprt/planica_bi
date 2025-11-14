package models

import "time"

// YandexCounter represents a Yandex.Metrica counter linked to a project
type YandexCounter struct {
	ID        uint      `gorm:"primaryKey"`
	ProjectID uint      `gorm:"not null;index"`
	CounterID int64     `gorm:"not null"`
	Name      string    // Optional name for the counter
	IsPrimary bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
