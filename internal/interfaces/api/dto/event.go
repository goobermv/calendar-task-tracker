package dto

import "time"

type CreateEventRequest struct {
	Title       string    `json:"title" binging:"required,min=1,max=255"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	EventType   string    `json:"event_type" binding:"omitempty,oneof=formal celebration important personal holiday casual business"`
}

type UpdateEventRequest struct {
	Title       *string    `json:"title" binging:"required,min=1,max=255"`
	Description *string    `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	EventType   *string    `json:"event_type" binding:"omitempty,oneof=formal celebration important personal holiday casual business"`
}

type EventResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	EventType   string    `json:"event_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
