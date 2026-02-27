package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/suprt/planica_bi/backend/internal/config"
	applogger "github.com/suprt/planica_bi/backend/internal/logger"
	"github.com/suprt/planica_bi/backend/internal/models"
)

// DB is the global database connection
var DB *gorm.DB

// Connect initializes database connection with retry logic
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase,
	)

	var err error
	var db *gorm.DB

	// Retry up to 30 times with 2 second delays
	for i := 0; i < 30; i++ {
		logMode := logger.Silent
		if cfg.AppDebug {
			logMode = logger.Info
		}
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logMode),
		})
		if err == nil {
			break
		}
		if applogger.Log != nil {
			applogger.Log.Warn("Database connection failed, retrying...",
				zap.Error(err),
				zap.Int("attempt", i+1),
			)
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after 30 attempts: %w", err)
	}

	DB = db

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if applogger.Log != nil {
		applogger.Log.Info("Database connection established")
	}
	return DB, nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// AutoMigrate runs database migrations using GORM
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	applogger.Log.Info("Starting GORM auto-migration...")

	// Run GORM auto migrations for models
	err := DB.AutoMigrate(
		&models.Project{},
		&models.User{},
		&models.UserProjectRole{},
		&models.YandexCounter{},
		&models.DirectAccount{},
		&models.DirectCampaign{},
		&models.Goal{},
		&models.MetricsMonthly{},
		&models.MetricsAgeMonthly{},
		&models.DirectCampaignMonthly{},
		&models.DirectTotalsMonthly{},
		&models.SEOQueriesMonthly{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	applogger.Log.Info("Database migrations completed")
	return nil
}

// SeedData inserts seed data into the database
func SeedData() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	applogger.Log.Info("Seeding database...")

	// Create admin user if not exists
	var admin models.User
	result := DB.Where("email = ?", "admin@test.ru").First(&admin)
	if result.Error == gorm.ErrRecordNotFound {
		admin = models.User{
			Name:     "Администратор",
			Email:    "admin@test.ru",
			Password: "$2a$10$Zw3HnTNxTLv6Jv/0Ld8YWuOxXzDuCmgyFa.1aXB.SQgzYFuZHYBym",
			Timezone: "Europe/Moscow",
			Language: "ru",
			IsActive: true,
		}
		if err := DB.Create(&admin).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
		applogger.Log.Info("Admin user created", zap.String("email", admin.Email))
	}

	// Create test project if not exists
	var project models.Project
	result = DB.Where("slug = ?", "test-project").First(&project)
	if result.Error == gorm.ErrRecordNotFound {
		project = models.Project{
			Name:      "Тестовый проект",
			Slug:      "test-project",
			PublicToken: "test-token-123",
			Timezone:  "Europe/Moscow",
			Currency:  "RUB",
			IsActive:  true,
		}
		if err := DB.Create(&project).Error; err != nil {
			return fmt.Errorf("failed to create test project: %w", err)
		}
		applogger.Log.Info("Test project created", zap.String("slug", project.Slug))

		// Assign admin role to admin user for the test project
		userProjectRole := models.UserProjectRole{
			UserID:    admin.ID,
			ProjectID: project.ID,
			Role:      "admin",
		}
		if err := DB.Create(&userProjectRole).Error; err != nil {
			return fmt.Errorf("failed to create user project role: %w", err)
		}
	}

	applogger.Log.Info("Database seeding completed")
	return nil
}

// SeedTestData inserts test data for the test project
func SeedTestData() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	applogger.Log.Info("Seeding test data...")

	// Get test project
	var project models.Project
	if err := DB.Where("slug = ?", "test-project").First(&project).Error; err != nil {
		return fmt.Errorf("test project not found: %w", err)
	}

	// Check if metrics already exist
	var count int64
	if err := DB.Model(&models.MetricsMonthly{}).Where("project_id = ?", project.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing metrics: %w", err)
	}
	if count > 0 {
		applogger.Log.Info("Test data already exists, skipping")
		return nil
	}

	// Seed metrics_monthly for the last 3 months
	months := []struct {
		year               int
		month              int
		visits             int
		users              int
		bounceRate         float64
		avgSessionDuration int
		conversions        int
	}{
		{year: 2025, month: 12, visits: 15420, users: 12350, bounceRate: 35.5, avgSessionDuration: 245, conversions: 1230},
		{year: 2026, month: 1, visits: 18200, users: 14500, bounceRate: 32.0, avgSessionDuration: 280, conversions: 1450},
		{year: 2026, month: 2, visits: 16800, users: 13200, bounceRate: 33.2, avgSessionDuration: 265, conversions: 1340},
	}

	for _, m := range months {
		metric := models.MetricsMonthly{
			ProjectID:             project.ID,
			Year:                  m.year,
			Month:                 m.month,
			Visits:                m.visits,
			Users:                 m.users,
			BounceRate:            m.bounceRate,
			AvgSessionDurationSec: m.avgSessionDuration,
			Conversions:           &m.conversions,
		}
		if err := DB.Create(&metric).Error; err != nil {
			return fmt.Errorf("failed to create metrics for %d-%02d: %w", m.year, m.month, err)
		}
	}

	// Seed direct_totals_monthly
	directMonths := []struct {
		year        int
		month       int
		impressions int
		clicks      int
		ctr         float64
		cpc         float64
		conversions int
		cpa         float64
		cost        float64
	}{
		{year: 2025, month: 12, impressions: 120000, clicks: 5200, ctr: 4.33, cpc: 42.50, conversions: 195, cpa: 1125.00, cost: 221000.00},
		{year: 2026, month: 1, impressions: 145000, clicks: 6100, ctr: 4.21, cpc: 38.90, conversions: 220, cpa: 985.00, cost: 237290.00},
		{year: 2026, month: 2, impressions: 130000, clicks: 5700, ctr: 4.38, cpc: 39.50, conversions: 215, cpa: 1070.00, cost: 225140.00},
	}

	for _, m := range directMonths {
		conversions := m.conversions
		cpa := m.cpa
		direct := models.DirectTotalsMonthly{
			ProjectID:   project.ID,
			Year:        m.year,
			Month:       m.month,
			Impressions: m.impressions,
			Clicks:      m.clicks,
			CTRPct:      m.ctr,
			CPC:         m.cpc,
			Conversions: &conversions,
			CPA:         &cpa,
			Cost:        m.cost,
		}
		if err := DB.Create(&direct).Error; err != nil {
			return fmt.Errorf("failed to create direct totals for %d-%02d: %w", m.year, m.month, err)
		}
	}

	applogger.Log.Info("Test data seeding completed")
	return nil
}
