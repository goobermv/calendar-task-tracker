package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in progress"
	TaskStatusCompleted  = "completed"
	TaskStatusCancelled  = "cancelled"
)

const (
	TaskPriorityUrgent = "urgent"
	TaskPriorityHigh   = "high"
	TaskPriorityMedium = "medium"
	TaskPriorityLow    = "low"
)

type Task struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Status      string
	Priority    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DueDate     time.Time
}
