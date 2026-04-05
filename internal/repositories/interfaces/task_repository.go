package repositories

import (
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	FindByID(id uuid.UUID) (*domain.Task, error)
	FindByTitle(title string) ([]*domain.Task, error)
	FindByDate(date time.Time) ([]*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id uuid.UUID) error
}
