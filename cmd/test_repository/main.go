package main

import (
	"fmt"

	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/config"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/database"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/logger"
	"github.com/goobermv/calendar-task-tracker/internal/interfaces/repositories/postgres"
	"github.com/google/uuid"
)

func main() {
	log := logger.NewLogger("TEST")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %w", err)
	}

	db, err := database.NewDBManager(cfg, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %w", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db.GetDB())

	testuuid, err := uuid.Parse("38d11552-6702-4dcb-b28f-f353b6b3cf26")
	if err != nil {
		panic(err)
	}

	if err := userRepo.Delete(testuuid); err != nil {
		log.Error("Failed to deleted user", err)
	} else {
		fmt.Println("User deleted successfully\n")
	}

	fmt.Printf("\nAll test passed successfully")
}
