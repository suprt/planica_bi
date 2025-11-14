package models

import "time"

// Goal represents a Yandex.Metrica goal
type Goal struct {
	ID           uint      `gorm:"primaryKey"`
	CounterID    uint      `gorm:"not null;index"`
	GoalID       int64     `gorm:"not null"`
	Name         string    // Optional goal name
	IsConversion bool      `gorm:"default:false"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
