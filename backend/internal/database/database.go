package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	applogger "gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
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
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
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

	if applogger.Log != nil {
		applogger.Log.Info("Database connection established")
	}
	return DB, nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

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

	if applogger.Log != nil {
		applogger.Log.Info("Database migrations completed")
	}
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
