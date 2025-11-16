package models

import "time"

// DirectAccount represents a Yandex.Direct account linked to a project
type DirectAccount struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"not null;index" json:"project_id"`
	ClientLogin string    `gorm:"type:varchar(255);not null" json:"client_login"`
	AccountName *string   `gorm:"type:varchar(255)" json:"account_name"` // Optional account name
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
