package models

// DirectAccount represents a Yandex.Direct account linked to a project
type DirectAccount struct {
	ID          uint   `gorm:"primaryKey"`
	ProjectID   uint   `gorm:"not null;index"`
	ClientLogin string `gorm:"not null"`
	AccountName string
	CreatedAt   int64
	UpdatedAt   int64
}
