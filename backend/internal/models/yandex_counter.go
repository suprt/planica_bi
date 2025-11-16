package models

import "time"

// YandexCounter represents a Yandex.Metrica counter linked to a project
type YandexCounter struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"not null;index" json:"project_id"`
	CounterID int64     `gorm:"not null" json:"counter_id"`
	Name      string    `json:"name"` // Optional name for the counter
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
