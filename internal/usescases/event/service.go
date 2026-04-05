package event

import (
	"errors"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	repositories "github.com/goobermv/calendar-task-tracker/internal/repositories/interfaces"
	"github.com/google/uuid"
)

type Service struct {
	eventRepo repositories.EventRepository
	userRepo  repositories.UserRepository
}

func NewService(eventRepo repositories.EventRepository, userRepo repositories.UserRepository) *Service {
	return &Service{
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

type CreateEventRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	EventType   string    `json:"event_type"`
}

func (s *Service) CreateEvent(req CreateEventRequest) (*domain.Event, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if len(req.Title) > 255 {
		return nil, errors.New("title must be less than 255 characters")
	}

	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if !req.StartTime.IsZero() && req.StartTime.Before(time.Now()) {
		return nil, errors.New("start time cannot be in the past")
	}
	if !req.EndTime.IsZero() && req.EndTime.Before(time.Now()) {
		return nil, errors.New("start time cannot be in the past")
	}
	if !req.StartTime.IsZero() && req.StartTime.After(req.EndTime) {
		return nil, errors.New("start time cannot be after end time")
	}
	if !req.EndTime.IsZero() && req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time cannot be before start time")
	}

	now := time.Now()
	event := &domain.Event{
		ID:          uuid.New(),
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		EventType:   req.EventType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err = s.eventRepo.Create(event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

type UpdateEventRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	EventType   *string    `json:"event_type"`
}

func (s *Service) UpdateEvent(eventID, userID uuid.UUID, req UpdateEventRequest) (*domain.Event, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find event by ID: %w", err)
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	if event.UserID != userID {
		return nil, fmt.Errorf("you do not have the permissions to update this event")
	}

	if req.Title != nil {
		if len(*req.Title) > 255 {
			return nil, errors.New("title cannot be longer than 255 characters")
		}
		if len(*req.Title) == 0 {
			return nil, errors.New("title cannot be empty")
		}
		event.Title = *req.Title
	}

	if req.Description != nil {
		event.Description = *req.Description
	}

	if req.StartTime != nil {
		if !req.StartTime.IsZero() && req.StartTime.Before(time.Now()) {
			return nil, errors.New("start time cannot be in the past")
		}
		if !req.StartTime.IsZero() && req.StartTime.After(*req.EndTime) {
			return nil, errors.New("start time cannot be after end time")
		}
		event.StartTime = *req.StartTime
	}

	if req.EndTime != nil {
		if !req.EndTime.IsZero() && req.EndTime.Before(time.Now()) {
			return nil, errors.New("start time cannot be in the past")
		}
		if !req.EndTime.IsZero() && req.EndTime.Before(*req.StartTime) {
			return nil, errors.New("start time cannot be after end time")
		}
		event.EndTime = *req.EndTime
	}

	if req.EventType != nil {
		validEventType := map[string]bool{
			domain.EventTypeFormal:      true,
			domain.EventTypeCelebration: true,
			domain.EventTypeImportant:   true,
			domain.EventTypePersonal:    true,
			domain.EventTypeHoliday:     true,
			domain.EventTypeCasual:      true,
			domain.EventTypeBusiness:    true,
		}
		if !validEventType[*req.EventType] {
			return nil, errors.New("invalid event type value")
		}
		event.EventType = *req.EventType
	}

	event.UpdatedAt = time.Now()

	if err := s.eventRepo.Update(event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return event, nil
}

func (s *Service) DeleteEvent(eventID, userID uuid.UUID) error {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return fmt.Errorf("failed to find event by ID: %w", err)
	}

	if event == nil {
		return errors.New("event not found")
	}

	if event.UserID != userID {
		return errors.New("you do not have the permissions to delete this task")
	}

	if err := s.eventRepo.Delete(eventID); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func (s *Service) GetEventByID(eventID, userID uuid.UUID) (*domain.Event, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find event by ID: %w", err)
	}

	if event == nil {
		return nil, errors.New("event not found")
	}

	if event.UserID != userID {
		return nil, errors.New("you do not have the permissions to delete this task")
	}

	return event, nil
}

// implement GetUserEvents function
// implement more methods that might come in handy and think about updating the already existing methods to be better
