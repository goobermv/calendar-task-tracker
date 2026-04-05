package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	ServerPort string

	User     string
	Port     int
	Host     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() (*DBConfig, error) {
	if err := LoadEnvFile(); err != nil {
		panic("reading db info from env error" + err.Error())
	}

	config := &DBConfig{
		ServerPort: GetEnv("SERVER_PORT", "8080"),
		User:       GetEnv("DB_USER", "postgres"),
		Port:       GetEnvAsInt("DB_PORT", 5432),
		Host:       GetEnv("DB_HOST", "localhost"),
		Password:   GetEnv("DB_PASSWORD", "4601"),
		DBName:     GetEnv("DB_NAME", "calendar_task_tracker_db"),
		SSLMode:    GetEnv("DB_SSLMODE", "disable"),
	}

	return config, nil
}

func LoadEnvFile() error {
	if err := godotenv.Load(".env.dbinfo"); err == nil {
		return nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get working directory")
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			envPath := filepath.Join(dir, ".env.dbinfo")
			if err := godotenv.Load(envPath); err == nil {
				return nil
			}
			return fmt.Errorf("no .env.dbinfo found at %s", envPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return fmt.Errorf("could not find project root or .env.dbinfo file")
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf(
		"user=%s port=%d host=%s password=%s dbname=%s sslmode=%s",
		c.User, c.Port, c.Host, c.Password, c.DBName, c.SSLMode,
	)
}

func (c *DBConfig) GetDBConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}
