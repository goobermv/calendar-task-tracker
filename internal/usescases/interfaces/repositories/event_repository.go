package repositories

import (
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/google/uuid"
)

type EventRepository interface {
	Create(event *domain.Event) error
	FindByID(id uuid.UUID) (*domain.Event, error)
	FindByTitle(title string) ([]*domain.Event, error)
	FindByDate(start, end time.Time) ([]*domain.Event, error)
	Update(task *domain.Event) error
	Delete(id uuid.UUID) error
}
