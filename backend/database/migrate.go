package database

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs database migrations
func RunMigrations(dbURL string) error {
	// Get absolute path to migrations directory
	migrationsPath, err := filepath.Abs("database/migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Run all up migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// RollbackMigrations rolls back database migrations
func RollbackMigrations(dbURL string, steps int) error {
	m, err := migrate.New(
		"file://database/migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Rollback specified number of migrations
	for i := 0; i < steps; i++ {
		if err := m.Steps(-1); err != nil {
			if err == migrate.ErrNoChange {
				break
			}
			return fmt.Errorf("failed to rollback migration: %w", err)
		}
	}

	return nil
}
