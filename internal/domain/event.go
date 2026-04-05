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
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	EventType   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
