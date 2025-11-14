package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
)

// DB is the global database connection
var DB *gorm.DB

// Connect initializes database connection
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established")
	return DB, nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	err := DB.AutoMigrate(
		&models.Project{},
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

	log.Println("Database migrations completed")
	return nil
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
