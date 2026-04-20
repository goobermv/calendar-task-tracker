package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title" binging:"required,min=1,max=255"`
	Description string    `json:"description"`
	Priority    string    `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string    `json:"description"`
	Status      *string    `json:"status" binding:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`
}

type TaskResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DueDate     time.Time `json:"due_datedom"`
}
