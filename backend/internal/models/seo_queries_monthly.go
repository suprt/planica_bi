package models

import "time"

// SEOQueriesMonthly represents monthly SEO query positions
type SEOQueriesMonthly struct {
	ID        uint      `gorm:"primaryKey"`
	ProjectID uint      `gorm:"not null;index"`
	Year      int       `gorm:"not null;index"`
	Month     int       `gorm:"not null;index"`
	Query     string    `gorm:"not null"`
	Position  int       `gorm:"not null"`
	URL       *string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

