package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventTypeFormal      = "formal"
	EventTypeCelebration = "celebration"
	EventTypeImportant   = "important"
	EventTypePersonal    = "personal"
	EventTypeHoliday     = "holiday"
	EventTypeCasual      = "casual"
	EventTypeBusiness    = "business"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	EventType   string    `json:"event_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
