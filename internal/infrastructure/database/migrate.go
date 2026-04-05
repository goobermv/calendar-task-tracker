package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func NewMigrator(databaseURL, migrationsPath string) (*Migrator, error) {
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{migrate: m}, nil
}

func (m *Migrator) Up() error {
	if err := m.migrate.Up(); err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migration: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

func (m *Migrator) Down() error {
	if err := m.migrate.Down(); err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rolled back successfully")
	return nil
}

func (m *Migrator) Steps(n int) error {
	if err := m.migrate.Steps(n); err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply %d migrations: %w", n, err)
	}

	log.Printf("Applied %d migrations successfully\n", n)
	return nil
}

func (m *Migrator) Force(version int) error {
	if err := m.migrate.Force(version); err != nil {
		return fmt.Errorf("failed to force version: %w", err)
	}

	return nil
}

func (m *Migrator) Version() (uint, bool, error) {
	return m.migrate.Version()
}

func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return sourceErr
	}
	if dbErr != nil {
		return dbErr
	}

	return nil
}
