package models

// SEOQueriesMonthly represents monthly SEO query positions
type SEOQueriesMonthly struct {
	ID        uint   `gorm:"primaryKey"`
	ProjectID uint   `gorm:"not null;index"`
	Year      int    `gorm:"not null;index"`
	Month     int    `gorm:"not null;index"`
	Query     string
	Position  int
	URL       *string
	CreatedAt int64
}

