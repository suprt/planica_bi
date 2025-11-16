package models

import "time"

// SEOQueriesMonthly represents monthly SEO query positions
type SEOQueriesMonthly struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"not null;index" json:"project_id"`
	Year      int       `gorm:"not null;index" json:"year"`
	Month     int       `gorm:"not null;index" json:"month"`
	Query     string    `gorm:"not null" json:"query"`
	Position  int       `gorm:"not null" json:"position"`
	URL       *string   `json:"url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for SEOQueriesMonthly
func (SEOQueriesMonthly) TableName() string {
	return "seo_queries_monthly"
}

