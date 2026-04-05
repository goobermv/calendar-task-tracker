package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/auth"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/config"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/database"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/logger"
	"github.com/goobermv/calendar-task-tracker/internal/network/api"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/handlers"
	"github.com/goobermv/calendar-task-tracker/internal/repositories/postgres"
	eventUsecases "github.com/goobermv/calendar-task-tracker/internal/usescases/event"
	taskUsecases "github.com/goobermv/calendar-task-tracker/internal/usescases/task"
	userUsecases "github.com/goobermv/calendar-task-tracker/internal/usescases/user"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger := logger.NewLogger("API")

	db, err := database.NewDBManager(cfg, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %w", err)
	}
	defer db.Close()

	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJWTService("secret_key", 1*time.Hour)

	userRepo := postgres.NewUserRepository(db.GetDB())
	taskRepo := postgres.NewTaskRepository(db.GetDB())
	eventRepo := postgres.NewEventRepository(db.GetDB())

	userService := userUsecases.NewService(userRepo, passwordService, jwtService)
	taskService := taskUsecases.NewService(taskRepo, userRepo)
	eventService := eventUsecases.NewService(eventRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)
	eventHandler := handlers.NewEventHandler(eventService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	api.SetupRoutes(router, userHandler, taskHandler, eventHandler, jwtService)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		appLogger.Infof("Server starting", "port", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatalf("Server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown: %w", err)
	}

	appLogger.Info("Server exited")
}
