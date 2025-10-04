package db

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// Migrate runs database migrations using Goose
func Migrate(db *sql.DB, migrationsDir string) error {
	// Run migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	return nil
}

// MigrateDown rolls back migrations
func MigrateDown(db *sql.DB, migrationsDir string) error {
	if err := goose.Down(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}
	
	return nil
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(db *sql.DB, migrationsDir string) (int64, error) {
	version, err := goose.GetDBVersion(db)
	if err != nil {
		return 0, fmt.Errorf("failed to get migration status: %w", err)
	}
	
	return version, nil
}