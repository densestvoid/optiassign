package db

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// Migrate runs database migrations using Goose
func Migrate(db *sql.DB, migrationsDir string) error {
	// Set the migrations directory
	goose.SetBaseDir(migrationsDir)
	
	// Run migrations
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	return nil
}

// MigrateDown rolls back migrations
func MigrateDown(db *sql.DB, migrationsDir string) error {
	goose.SetBaseDir(migrationsDir)
	
	if err := goose.Down(db, "."); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}
	
	return nil
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(db *sql.DB, migrationsDir string) ([]goose.MigrationRecord, error) {
	goose.SetBaseDir(migrationsDir)
	
	records, err := goose.GetDBVersion(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration status: %w", err)
	}
	
	return []goose.MigrationRecord{{Version: records}}, nil
}