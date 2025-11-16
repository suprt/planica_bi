package models

import "time"

// User represents a system user
type User struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Name            string     `gorm:"not null" json:"name"`
	Email           string     `gorm:"unique;not null;index" json:"email"`
	Password        string     `gorm:"not null" json:"-"` // Hashed password, exclude from JSON
	EmailVerifiedAt *time.Time `gorm:"type:timestamp NULL" json:"email_verified_at"`
	Timezone        string     `gorm:"default:Europe/Moscow" json:"timezone"`
	Language        string     `gorm:"type:enum('ru','en');default:'ru'" json:"language"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt     *time.Time `gorm:"type:timestamp NULL" json:"last_login_at"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	ProjectRoles []UserProjectRole `gorm:"foreignKey:UserID" json:"project_roles,omitempty"`
}

// UserProjectRole represents user role in a project (pivot table)
type UserProjectRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_project" json:"user_id"`
	ProjectID uint      `gorm:"not null;uniqueIndex:idx_user_project" json:"project_id"`
	Role      string    `gorm:"type:enum('admin','manager','client');not null" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}
