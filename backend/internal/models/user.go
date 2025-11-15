package models

import "time"

// User represents a system user
type User struct {
	ID              uint       `gorm:"primaryKey"`
	Name            string     `gorm:"not null"`
	Email           string     `gorm:"unique;not null;index"`
	Password        string     `gorm:"not null"` // Hashed password
	EmailVerifiedAt *time.Time `gorm:"type:timestamp NULL"`
	Timezone        string     `gorm:"default:Europe/Moscow"`
	Language        string     `gorm:"type:enum('ru','en');default:'ru'"`
	IsActive        bool       `gorm:"default:true"`
	LastLoginAt     *time.Time `gorm:"type:timestamp NULL"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`

	// Relations
	ProjectRoles []UserProjectRole `gorm:"foreignKey:UserID"`
}

// UserProjectRole represents user role in a project (pivot table)
type UserProjectRole struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_project"`
	ProjectID uint      `gorm:"not null;uniqueIndex:idx_user_project"`
	Role      string    `gorm:"type:enum('admin','manager','client');not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relations
	User    User    `gorm:"foreignKey:UserID"`
	Project Project `gorm:"foreignKey:ProjectID"`
}
