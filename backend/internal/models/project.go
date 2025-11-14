package models

// Project represents a client project
type Project struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Slug      string `gorm:"unique;not null"`
	Timezone  string `gorm:"default:Europe/Moscow"`
	Currency  string `gorm:"default:RUB"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt int64
	UpdatedAt int64
}

