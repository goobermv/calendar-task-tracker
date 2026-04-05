package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	UserTypeRegular = "regular"
	UserTypeAdmin   = "admin"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash string
	UserType     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
