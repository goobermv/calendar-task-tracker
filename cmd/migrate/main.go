package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/config"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/database"
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: up, down, steps, or force")
	steps := flag.Int("steps", 0, "Number of steps to apply (only used with direction=steps)")
	version := flag.Int("version", 0, "Version to force (only used with direction=force)")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Connecting to database: %s@%s:%d/%s\n", cfg.User, cfg.Host, cfg.Port, cfg.DBName)

	databaseURL := cfg.GetDBConnString()
	if databaseURL == "" {
		log.Fatal("Database URL is empty")
	}

	migrator, err := database.NewMigrator(databaseURL, "file://migrations")
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer migrator.Close()

	switch *direction {
	case "up":
		fmt.Println("Running migrations up...")
		if err := migrator.Up(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations applied successfully!")

	case "down":
		fmt.Println("Running migrations down...")
		if err := migrator.Down(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations rolled back successfully!")

	case "steps":
		if *steps <= 0 {
			log.Fatal("Steps must be greater than 0 when using direction=steps")
		}
		fmt.Printf("Running %d migration steps...\n", *steps)
		if err := migrator.Steps(*steps); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Printf("Applied %d migration steps successfully!\n", *steps)

	case "force":
		if *version < 0 {
			log.Fatal("Version must be >= 0 when using direction=force")
		}
		fmt.Printf("Forcing migration version to %d...\n", *version)
		if err := migrator.Force(*version); err != nil {
			log.Fatalf("Force failed: %v", err)
		}
		fmt.Printf("Successfully forced version to %d\n", *version)

	default:
		log.Fatalf("Unknown direction: %s. Use up, down, steps, or force", *direction)
	}

	currentVersion, dirty, err := migrator.Version()
	if err != nil {
		fmt.Printf("Warning: Could not get migration version: %v\n", err)
	} else {
		fmt.Printf("Current migration version: %d (dirty: %v)\n", currentVersion, dirty)
	}
}
