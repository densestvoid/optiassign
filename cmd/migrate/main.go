package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"optiassign/config"
	"optiassign/db"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Parse command line flags
	var (
		command = flag.String("command", "up", "Migration command: up, down, status")
		dir     = flag.String("dir", "migrations", "Migrations directory")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	if err := db.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Execute migration command
	switch *command {
	case "up":
		if err := db.Migrate(db.GetDB(), *dir); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
		fmt.Println("Migrations completed successfully")
	case "down":
		if err := db.MigrateDown(db.GetDB(), *dir); err != nil {
			log.Fatal("Failed to rollback migrations:", err)
		}
		fmt.Println("Migrations rolled back successfully")
	case "status":
		version, err := db.GetMigrationStatus(db.GetDB(), *dir)
		if err != nil {
			log.Fatal("Failed to get migration status:", err)
		}
		fmt.Printf("Current migration version: %d\n", version)
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: up, down, status")
		os.Exit(1)
	}
}