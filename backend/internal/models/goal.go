package models

// Goal represents a Yandex.Metrica goal
type Goal struct {
	ID           uint  `gorm:"primaryKey"`
	CounterID    uint  `gorm:"not null;index"`
	GoalID       int64 `gorm:"not null"`
	Name         string
	IsConversion bool `gorm:"default:false"`
	CreatedAt    int64
	UpdatedAt    int64
}
