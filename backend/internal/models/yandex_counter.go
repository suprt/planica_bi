package models

// YandexCounter represents a Yandex.Metrica counter linked to a project
type YandexCounter struct {
	ID        uint   `gorm:"primaryKey"`
	ProjectID uint   `gorm:"not null;index"`
	CounterID int64  `gorm:"not null"`
	Name      string
	IsPrimary bool   `gorm:"default:false"`
	CreatedAt int64
	UpdatedAt int64
}

