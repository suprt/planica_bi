package models

import "time"

// Goal represents a Yandex.Metrica goal
type Goal struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CounterID    uint      `gorm:"not null;index" json:"counter_id"`
	GoalID       int64     `gorm:"not null" json:"goal_id"`
	Name         string    `json:"name"` // Optional goal name
	IsConversion bool      `gorm:"default:false" json:"is_conversion"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
