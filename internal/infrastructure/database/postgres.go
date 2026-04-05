package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/config"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/logger"
	_ "github.com/lib/pq"
)

type DBManager struct {
	db     *sql.DB
	config *config.DBConfig
	logger logger.Logger
}

func NewDBManager(cfg *config.DBConfig, log logger.Logger) (*DBManager, error) {
	log.Info("starting connection to database")

	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("successfuly connected to database")

	return &DBManager{
		db:     db,
		config: cfg,
		logger: log,
	}, nil
}

func (db *DBManager) Close() error {
	if db.db != nil {
		if err := db.db.Close(); err != nil {
			db.logger.Error("closing database connection error", err)
			return err
		}
		db.logger.Success("database connection closed")
	}
	return nil
}

func (db *DBManager) GetDB() *sql.DB {
	return db.db
}
