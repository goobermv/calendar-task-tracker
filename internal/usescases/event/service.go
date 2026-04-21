package event

import (
	"errors"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/infrastructure/logger"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	repositories "github.com/goobermv/calendar-task-tracker/internal/repositories/interfaces"
	"github.com/google/uuid"
)

type Service struct {
	eventRepo repositories.EventRepository
	userRepo  repositories.UserRepository
	logger    logger.Logger
}

func NewService(eventRepo repositories.EventRepository, userRepo repositories.UserRepository, logger logger.Logger) *Service {
	return &Service{
		eventRepo: eventRepo,
		userRepo:  userRepo,
		logger:    logger,
	}
}

func (s *Service) CreateEvent(req dto.CreateEventRequest) (*domain.Event, error) {
	s.logger.Infof("Creating event for user: %s", req.UserID)

	if req.Title == "" {
		s.logger.Warnf("Event creation failed: title is required for user: %s", req.UserID)
		return nil, errors.New("title is required")
	}
	if len(req.Title) > 255 {
		s.logger.Warnf("Event creation failed: title too long (%d chars) for user: %s", len(req.Title), req.UserID)
		return nil, errors.New("title must be less than 255 characters")
	}

	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		s.logger.Errorf(err, "Failed to verify user: %s", req.UserID)
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}
	if user == nil {
		s.logger.Warnf("Event creation failed: user not found - %s", req.UserID)
		return nil, domain.ErrUserNotFound
	}

	if !req.StartTime.IsZero() && req.StartTime.Before(time.Now()) {
		s.logger.Warnf("Event creation failed: start time in the past for user: %s", req.UserID)
		return nil, errors.New("start time cannot be in the past")
	}
	if !req.EndTime.IsZero() && req.EndTime.Before(time.Now()) {
		s.logger.Warnf("Event creation failed: end time in the past for user: %s", req.UserID)
		return nil, errors.New("start time cannot be in the past")
	}
	if !req.StartTime.IsZero() && req.StartTime.After(req.EndTime) {
		s.logger.Warnf("Event creation failed: start time after end time for user: %s", req.UserID)
		return nil, errors.New("start time cannot be after end time")
	}
	if !req.EndTime.IsZero() && req.EndTime.Before(req.StartTime) {
		s.logger.Warnf("Event creation failed: end time before start time for user: %s", req.UserID)
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
		s.logger.Errorf(err, "Failed to create event in database for user: %s", req.UserID)
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	s.logger.Successf("Event created successfully - ID: %s, Title: %s, User: %s", event.ID, event.Title, req.UserID)
	return event, nil
}

func (s *Service) UpdateEvent(eventID, userID uuid.UUID, req dto.UpdateEventRequest) (*domain.Event, error) {
	s.logger.Debugf("Updating event: %s for user: %s", eventID, userID)

	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find event: %s", eventID)
		return nil, fmt.Errorf("failed to find event by ID: %w", err)
	}
	if event == nil {
		s.logger.Warnf("Event not found: %s", eventID)
		return nil, domain.ErrEventNotFound
	}

	if event.UserID != userID {
		s.logger.Warnf("Permission denied: user %s attempted to update event %s owned by %s", userID, eventID, event.UserID)
		return nil, domain.ErrNoPermissionUpdateEvent
	}

	if req.Title != nil {
		if len(*req.Title) > 255 {
			s.logger.Warnf("Event update failed: title too long for event: %s", eventID)
			return nil, errors.New("title cannot be longer than 255 characters")
		}
		if len(*req.Title) == 0 {
			s.logger.Warnf("Event update failed: empty title for event: %s", eventID)
			return nil, errors.New("title cannot be empty")
		}
		event.Title = *req.Title
		s.logger.Debugf("Event title updated to: %s", *req.Title)
	}

	if req.Description != nil {
		event.Description = *req.Description
		s.logger.Debugf("Event description updated for event: %s", eventID)
	}

	if req.StartTime != nil {
		if !req.StartTime.IsZero() && req.StartTime.Before(time.Now()) {
			s.logger.Warnf("Event update failed: start time in the past for event: %s", eventID)
			return nil, errors.New("start time cannot be in the past")
		}
		if !req.StartTime.IsZero() && req.StartTime.After(*req.EndTime) {
			s.logger.Warnf("Event update failed: start time after end time for event: %s", eventID)
			return nil, errors.New("start time cannot be after end time")
		}
		event.StartTime = *req.StartTime
		s.logger.Debugf("Event start time updated to: %v", *req.StartTime)
	}

	if req.EndTime != nil {
		if !req.EndTime.IsZero() && req.EndTime.Before(time.Now()) {
			s.logger.Warnf("Event update failed: end time in the past for event: %s", eventID)
			return nil, errors.New("start time cannot be in the past")
		}
		if !req.EndTime.IsZero() && req.EndTime.Before(*req.StartTime) {
			s.logger.Warnf("Event update failed: end time before start time for event: %s", eventID)
			return nil, errors.New("start time cannot be after end time")
		}
		event.EndTime = *req.EndTime
		s.logger.Debugf("Event end time updated to: %v", *req.EndTime)
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
			s.logger.Warnf("Event update failed: invalid event type: %s", *req.EventType)
			return nil, errors.New("invalid event type value")
		}
		event.EventType = *req.EventType
		s.logger.Debugf("Event type updated to: %s", *req.EventType)
	}

	event.UpdatedAt = time.Now()

	if err := s.eventRepo.Update(event); err != nil {
		s.logger.Errorf(err, "Failed to update event in database: %s", eventID)
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	s.logger.Successf("Event updated successfully - ID: %s, User: %s", eventID, userID)
	return event, nil
}

func (s *Service) DeleteEvent(eventID, userID uuid.UUID) error {
	s.logger.Infof("Deleting event: %s for user: %s", eventID, userID)

	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find event: %s", eventID)
		return fmt.Errorf("failed to find event by ID: %w", err)
	}

	if event == nil {
		s.logger.Warnf("Event not found for deletion: %s", eventID)
		return domain.ErrEventNotFound
	}

	if event.UserID != userID {
		s.logger.Warnf("Permission denied: user %s attempted to delete event %s owned by %s", userID, eventID, event.UserID)
		return domain.ErrNoPermissionDeleteEvent
	}

	if err := s.eventRepo.Delete(eventID); err != nil {
		s.logger.Errorf(err, "Failed to delete event from database: %s", eventID)
		return fmt.Errorf("failed to delete event: %w", err)
	}

	s.logger.Successf("Event deleted successfully - ID: %s, Title: %s", eventID, event.Title)
	return nil
}

func (s *Service) GetEventByID(eventID, userID uuid.UUID) (*domain.Event, error) {
	s.logger.Debugf("Fetching event: %s for user: %s", eventID, userID)

	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find event: %s", eventID)
		return nil, fmt.Errorf("failed to find event by ID: %w", err)
	}

	if event == nil {
		s.logger.Warnf("Event not found: %s", eventID)
		return nil, domain.ErrEventNotFound
	}

	if event.UserID != userID {
		s.logger.Warnf("Permission denied: user %s attempted to view event %s owned by %s", userID, eventID, event.UserID)
		return nil, domain.ErrNoPermissionViewEvent
	}

	s.logger.Debugf("Event retrieved successfully: %s", eventID)
	return event, nil
}

func (s *Service) GetUserEvents(userID uuid.UUID) ([]*domain.Event, error) {
	s.logger.Debugf("Fetching all events for user: %s", userID)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Errorf(err, "Failed to verify user: %s", userID)
		return nil, domain.ErrValidation
	}
	if user == nil {
		s.logger.Warnf("User not found for event retrieval: %s", userID)
		return nil, domain.ErrUserNotFound
	}

	events, err := s.eventRepo.FindByUserID(userID)
	if err != nil {
		s.logger.Errorf(err, "Failed to get events for user: %s", userID)
		return nil, fmt.Errorf("failed to get user tasks: %w", err)
	}

	s.logger.Debugf("Retrieved %d events for user: %s", len(events), userID)
	return events, nil
}

func (s *Service) AdminGetAllEvents() ([]*domain.Event, error) {
	s.logger.Infof("Admin: Fetching all events from all users")

	events, err := s.eventRepo.FindAll()
	if err != nil {
		s.logger.Errorf(err, "Admin: Failed to get all events")
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}

	s.logger.Infof("Admin: Retrieved %d total events", len(events))
	return events, nil
}

func (s *Service) AdminUpdateEvent(eventID uuid.UUID, req dto.AdminUpdateEventRequest) (*domain.Event, error) {
	s.logger.Infof("Admin: Updating event: %s", eventID)

	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		s.logger.Errorf(err, "Admin: Failed to find event: %s", eventID)
		return nil, fmt.Errorf("failed to find event by ID: %w", err)
	}
	if event == nil {
		s.logger.Warnf("Admin: Event not found for update: %s", eventID)
		return nil, domain.ErrEventNotFound
	}

	if req.Title != nil {
		if len(*req.Title) > 255 {
			s.logger.Warnf("Admin: Event update failed: title too long for event: %s", eventID)
			return nil, errors.New("title cannot be longer than 255 characters")
		}
		if len(*req.Title) == 0 {
			s.logger.Warnf("Admin: Event update failed: empty title for event: %s", eventID)
			return nil, errors.New("title cannot be empty")
		}
		event.Title = *req.Title
		s.logger.Debugf("Admin: Event title updated to: %s", *req.Title)
	}

	if req.Description != nil {
		event.Description = *req.Description
		s.logger.Debugf("Admin: Event description updated for event: %s", eventID)
	}

	if req.StartTime != nil {
		if !req.StartTime.IsZero() && req.StartTime.Before(time.Now()) {
			s.logger.Warnf("Admin: Event update failed: start time in the past for event: %s", eventID)
			return nil, errors.New("start time cannot be in the past")
		}
		if req.EndTime != nil && !req.StartTime.IsZero() && req.StartTime.After(*req.EndTime) {
			s.logger.Warnf("Admin: Event update failed: start time after end time for event: %s", eventID)
			return nil, errors.New("start time cannot be after end time")
		}
		if req.EndTime == nil && !req.StartTime.IsZero() && req.StartTime.After(event.EndTime) {
			s.logger.Warnf("Admin: Event update failed: start time after existing end time for event: %s", eventID)
			return nil, errors.New("start time cannot be after end time")
		}
		event.StartTime = *req.StartTime
		s.logger.Debugf("Admin: Event start time updated to: %v", *req.StartTime)
	}

	if req.EndTime != nil {
		if !req.EndTime.IsZero() && req.EndTime.Before(time.Now()) {
			s.logger.Warnf("Admin: Event update failed: end time in the past for event: %s", eventID)
			return nil, errors.New("end time cannot be in the past")
		}
		startTime := event.StartTime
		if req.StartTime != nil {
			startTime = *req.StartTime
		}
		if !req.EndTime.IsZero() && req.EndTime.Before(startTime) {
			s.logger.Warnf("Admin: Event update failed: end time before start time for event: %s", eventID)
			return nil, errors.New("end time cannot be before start time")
		}
		event.EndTime = *req.EndTime
		s.logger.Debugf("Admin: Event end time updated to: %v", *req.EndTime)
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
			s.logger.Warnf("Admin: Event update failed: invalid event type: %s", *req.EventType)
			return nil, errors.New("invalid event type value")
		}
		event.EventType = *req.EventType
		s.logger.Debugf("Admin: Event type updated to: %s", *req.EventType)
	}

	event.UpdatedAt = time.Now()

	if err := s.eventRepo.Update(event); err != nil {
		s.logger.Errorf(err, "Admin: Failed to update event in database: %s", eventID)
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	s.logger.Successf("Admin: Event updated successfully - ID: %s, User: %s", eventID, event.UserID)
	return event, nil
}

func (s *Service) AdminDeleteEvent(eventID uuid.UUID) error {
	s.logger.Infof("Admin: Deleting event: %s", eventID)

	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		s.logger.Errorf(err, "Admin: Failed to find event: %s", eventID)
		return fmt.Errorf("failed to find event by ID: %w", err)
	}
	if event == nil {
		s.logger.Warnf("Admin: Event not found for deletion: %s", eventID)
		return domain.ErrEventNotFound
	}

	if err := s.eventRepo.Delete(eventID); err != nil {
		s.logger.Errorf(err, "Admin: Failed to delete event from database: %s", eventID)
		return fmt.Errorf("failed to delete event: %w", err)
	}

	s.logger.Successf("Admin: Event deleted successfully - ID: %s, Title: %s, Owner: %s", eventID, event.Title, event.UserID)
	return nil
}
