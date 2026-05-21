package task

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
	taskRepo repositories.TaskRepository
	userRepo repositories.UserRepository
	logger   logger.Logger
}

func NewService(taskRepo repositories.TaskRepository, userRepo repositories.UserRepository, logger logger.Logger) *Service {
	return &Service{
		taskRepo: taskRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *Service) CreateTask(req dto.CreateTaskRequest) (*domain.Task, error) {
	s.logger.Infof("Creating task for user: %s", req.UserID)

	if req.Title == "" {
		s.logger.Warnf("Task creation failed: title is required for user: %s", req.UserID)
		return nil, errors.New("title is required")
	}
	if len(req.Title) > 255 {
		s.logger.Warnf("Task creation failed: title too long (%d chars) for user: %s", len(req.Title), req.UserID)
		return nil, errors.New("title must be less than 255 characters")
	}

	if !req.DueDate.IsZero() && req.DueDate.Before(time.Now().Truncate(24*time.Hour)) {
		s.logger.Warnf("Task creation failed: due date in the past for user: %s", req.UserID)
		return nil, errors.New("due date cannot be in the past")
	}

	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		s.logger.Errorf(err, "Failed to verify user: %s", req.UserID)
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}
	if user == nil {
		s.logger.Warnf("Task creation failed: user not found - %s", req.UserID)
		return nil, domain.ErrUserNotFound
	}

	newID := uuid.New()
	now := time.Now()
	task := &domain.Task{
		ID:          newID,
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.TaskStatusPending,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.taskRepo.Create(task); err != nil {
		s.logger.Errorf(err, "Failed to create task in database for user: %s", req.UserID)
		return nil, domain.ErrFailedToCreateTask
	}

	task.ID = newID

	s.logger.Successf("Task created successfully - ID: %s, Title: %s, User: %s", task.ID, task.Title, req.UserID)
	return task, nil
}

func (s *Service) UpdateTask(taskID, userID uuid.UUID, req dto.UpdateTaskRequest) (*domain.Task, error) {
	s.logger.Debugf("Updating task: %s for user: %s", taskID, userID)

	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find task: %s", taskID)
		return nil, fmt.Errorf("failed to find task by ID: %w", err)
	}
	if task == nil {
		s.logger.Warnf("Task not found: %s", taskID)
		return nil, domain.ErrTaskNotFound
	}

	if task.UserID != userID {
		s.logger.Warnf("Permission denied: user %s attempted to update task %s owned by %s", userID, taskID, task.UserID)
		return nil, domain.ErrNoPermissionUpdateTask
	}

	if req.Title != nil {
		if len(*req.Title) > 255 {
			s.logger.Warnf("Task update failed: title too long for task: %s", taskID)
			return nil, errors.New("title cannot be longer than 255 characters")
		}
		if len(*req.Title) == 0 {
			s.logger.Warnf("Task update failed: empty title for task: %s", taskID)
			return nil, errors.New("title cannot be empty")
		}
		task.Title = *req.Title
		s.logger.Debugf("Task title updated to: %s", *req.Title)
	}

	if req.Description != nil {
		task.Description = *req.Description
		s.logger.Debugf("Task description updated for task: %s", taskID)
	}

	if req.DueDate != nil {
		if !req.DueDate.IsZero() && req.DueDate.Before(time.Now()) {
			s.logger.Warnf("Task update failed: due date in the past for task: %s", taskID)
			return nil, errors.New("due date cannot be in the past")
		}
		task.DueDate = *req.DueDate
		s.logger.Debugf("Task due date updated to: %v", *req.DueDate)
	}

	if req.Priority != nil {
		validPrioity := map[string]bool{
			domain.TaskPriorityUrgent: true,
			domain.TaskPriorityHigh:   true,
			domain.TaskPriorityMedium: true,
			domain.TaskPriorityLow:    true,
		}
		if !validPrioity[*req.Priority] {
			s.logger.Warnf("Task update failed: invalid priority value: %s", *req.Priority)
			return nil, errors.New("invalid priority value")
		}
		task.Priority = *req.Priority
		s.logger.Debugf("Task priority updated to: %s", *req.Priority)
	}

	if req.Status != nil {
		statusVal := *req.Status

		validStatus := map[string]bool{
			domain.TaskStatusPending:    true,
			domain.TaskStatusInProgress: true,
			domain.TaskStatusCompleted:  true,
			domain.TaskStatusCancelled:  true,
		}

		if !validStatus[statusVal] {
			s.logger.Warnf("Task update failed: invalid status value: %s", statusVal)
			return nil, fmt.Errorf("invalid status value: %s", statusVal)
		}

		task.Status = statusVal
		s.logger.Debugf("Task status updated to: %s", statusVal)
	}

	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(task); err != nil {
		s.logger.Errorf(err, "Failed to update task in database: %s", taskID)
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	s.logger.Successf("Task updated successfully - ID: %s, User: %s", taskID, userID)
	return task, nil
}

func (s *Service) DeleteTask(taskID, userID uuid.UUID) error {
	s.logger.Infof("Deleting task: %s for user: %s", taskID, userID)

	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find task: %s", taskID)
		return fmt.Errorf("failed to find task by ID: %w", err)
	}
	if task == nil {
		s.logger.Warnf("Task not found for deletion: %s", taskID)
		return domain.ErrTaskNotFound
	}

	if task.UserID != userID {
		s.logger.Warnf("Permission denied: user %s attempted to delete task %s owned by %s", userID, taskID, task.UserID)
		return domain.ErrNoPermissionDeleteTask
	}

	if err := s.taskRepo.Delete(taskID); err != nil {
		s.logger.Errorf(err, "Failed to delete task from database: %s", taskID)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	s.logger.Successf("Task deleted successfully - ID: %s, Title: %s", taskID, task.Title)
	return nil
}

func (s *Service) GetTaskByID(taskID, userID uuid.UUID) (*domain.Task, error) {
	s.logger.Debugf("Fetching task: %s for user: %s", taskID, userID)

	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		s.logger.Errorf(err, "Failed to find task: %s", taskID)
		return nil, fmt.Errorf("failed to find task by ID: %w", err)
	}
	if task == nil {
		s.logger.Warnf("Task not found: %s", taskID)
		return nil, domain.ErrTaskNotFound
	}

	if task.UserID != userID {
		s.logger.Warnf("Permission denied: user %s attempted to view task %s owned by %s", userID, taskID, task.UserID)
		return nil, domain.ErrNoPermissionViewTask
	}

	s.logger.Debugf("Task retrieved successfully: %s", taskID)
	return task, nil
}

func (s *Service) GetUserTasks(userID uuid.UUID, page, limit int) ([]*domain.Task, int, error) {
	s.logger.Debugf("Fetching all tasks for user: %s (page: %d, limit: %d)", userID, page, limit)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Errorf(err, "Failed to verify user: %s", userID)
		return nil, 0, domain.ErrValidation
	}
	if user == nil {
		s.logger.Warnf("User not found for task retrieval: %s", userID)
		return nil, 0, domain.ErrUserNotFound
	}

	offset := (page - 1) * limit

	tasks, totalCount, err := s.taskRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		s.logger.Errorf(err, "Failed to get tasks for user: %s", userID)
		return nil, 0, fmt.Errorf("failed to get user tasks: %w", err)
	}

	return tasks, totalCount, nil
}

func (s *Service) AdminGetAllTasks(page, limit int) ([]*domain.Task, int, error) {
	s.logger.Infof("Admin: Fetching tasks with pagination (page: %d, limit: %d)", page, limit)

	offset := (page - 1) * limit

	tasks, totalCount, err := s.taskRepo.FindAll(limit, offset)
	if err != nil {
		s.logger.Errorf(err, "Admin: Failed to get all tasks")
		return nil, 0, fmt.Errorf("failed to get all tasks: %w", err)
	}

	s.logger.Infof("Admin: Retrieved %d tasks out of %d total", len(tasks), totalCount)
	return tasks, totalCount, nil
}

func (s *Service) AdminUpdateTask(taskID uuid.UUID, req dto.AdminUpdateTaskRequest) (*domain.Task, error) {
	s.logger.Infof("Admin: Updating task: %s", taskID)

	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		s.logger.Errorf(err, "Admin: Failed to find task: %s", taskID)
		return nil, fmt.Errorf("failed to find task by ID: %w", err)
	}
	if task == nil {
		s.logger.Warnf("Admin: Task not found for update: %s", taskID)
		return nil, domain.ErrTaskNotFound
	}

	if req.Title != nil {
		if len(*req.Title) > 255 {
			s.logger.Warnf("Admin: Task update failed: title too long for task: %s", taskID)
			return nil, errors.New("title cannot be longer than 255 characters")
		}
		if len(*req.Title) == 0 {
			s.logger.Warnf("Admin: Task update failed: empty title for task: %s", taskID)
			return nil, errors.New("title cannot be empty")
		}
		task.Title = *req.Title
		s.logger.Debugf("Admin: Task title updated to: %s", *req.Title)
	}

	if req.Description != nil {
		task.Description = *req.Description
		s.logger.Debugf("Admin: Task description updated for task: %s", taskID)
	}

	if req.DueDate != nil {
		if !req.DueDate.IsZero() && req.DueDate.Before(time.Now()) {
			s.logger.Warnf("Admin: Task update failed: due date in the past for task: %s", taskID)
			return nil, errors.New("due date cannot be in the past")
		}
		task.DueDate = *req.DueDate
		s.logger.Debugf("Admin: Task due date updated to: %v", *req.DueDate)
	}

	if req.Priority != nil {
		validPrioity := map[string]bool{
			domain.TaskPriorityUrgent: true,
			domain.TaskPriorityHigh:   true,
			domain.TaskPriorityMedium: true,
			domain.TaskPriorityLow:    true,
		}
		if !validPrioity[*req.Priority] {
			s.logger.Warnf("Admin: Task update failed: invalid priority value: %s", *req.Priority)
			return nil, errors.New("invalid priority value")
		}
		task.Priority = *req.Priority
		s.logger.Debugf("Admin: Task priority updated to: %s", *req.Priority)
	}

	if req.Status != nil {
		statusVal := *req.Status

		validStatus := map[string]bool{
			domain.TaskStatusPending:    true,
			domain.TaskStatusInProgress: true,
			domain.TaskStatusCompleted:  true,
			domain.TaskStatusCancelled:  true,
		}

		if !validStatus[statusVal] {
			s.logger.Warnf("Task update failed: invalid status value: %s", statusVal)
			return nil, fmt.Errorf("invalid status value: %s", statusVal)
		}

		task.Status = statusVal
		s.logger.Debugf("Task status updated to: %s", statusVal)
	}

	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(task); err != nil {
		s.logger.Errorf(err, "Admin: Failed to update task in database: %s", taskID)
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	s.logger.Successf("Admin: Task updated successfully - ID: %s, User: %s", taskID, task.UserID)
	return task, nil
}

func (s *Service) AdminDeleteTask(taskID uuid.UUID) error {
	s.logger.Infof("Admin: Deleting task: %s", taskID)

	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		s.logger.Errorf(err, "Admin: Failed to find task: %s", taskID)
		return fmt.Errorf("failed to find task by ID: %w", err)
	}
	if task == nil {
		s.logger.Warnf("Admin: Task not found for deletion: %s", taskID)
		return domain.ErrTaskNotFound
	}

	if err := s.taskRepo.Delete(taskID); err != nil {
		s.logger.Errorf(err, "Admin: Failed to delete task from database: %s", taskID)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	s.logger.Successf("Admin: Task deleted successfully - ID: %s, Title: %s, Owner: %s", taskID, task.Title, task.UserID)
	return nil
}
