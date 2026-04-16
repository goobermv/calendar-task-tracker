package task

import (
	"errors"
	"fmt"
	"time"

	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	repositories "github.com/goobermv/calendar-task-tracker/internal/repositories/interfaces"

	"github.com/google/uuid"
)

type Service struct {
	taskRepo repositories.TaskRepository
	userRepo repositories.UserRepository
}

func NewService(taskRepo repositories.TaskRepository, userRepo repositories.UserRepository) *Service {
	return &Service{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (s *Service) CreateTask(req dto.CreateTaskRequest) (*domain.Task, error) {
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
		return nil, domain.ErrUserNotFound
	}

	now := time.Now()
	task := &domain.Task{
		ID:          uuid.New(),
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
		return nil, domain.ErrFailedToCreateTask
	}

	return task, nil
}

func (s *Service) UpdateTask(taskID, userID uuid.UUID, req dto.UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task by ID: %w", err)
	}
	if task == nil {
		return nil, domain.ErrTaskNotFound
	}

	if task.UserID != userID {
		return nil, domain.ErrNoPermissionUpdateTask
	}

	if req.Title != nil {
		if len(*req.Title) > 255 {
			return nil, errors.New("title cannot be longer than 255 characters")
		}
		if len(*req.Title) == 0 {
			return nil, errors.New("title cannot be empty")
		}
		task.Title = *req.Title
	}

	if req.Description != nil {
		task.Description = *req.Description
	}

	if req.DueDate != nil {
		if !req.DueDate.IsZero() && req.DueDate.Before(time.Now()) {
			return nil, errors.New("due date cannot be in the past")
		}
		task.DueDate = *req.DueDate
	}

	if req.Priority != nil {
		validPrioity := map[string]bool{
			domain.TaskPriorityUrgent: true,
			domain.TaskPriorityHigh:   true,
			domain.TaskPriorityMedium: true,
			domain.TaskPriorityLow:    true,
		}
		if !validPrioity[*req.Priority] {
			return nil, errors.New("invalid priority value")
		}
		task.Priority = *req.Priority
	}

	if req.Status != nil {
		validStatus := map[string]bool{
			domain.TaskStatusPending:    true,
			domain.TaskStatusInProgress: true,
			domain.TaskStatusCompleted:  true,
			domain.TaskStatusCancelled:  true,
		}
		if !validStatus[*req.Status] {
			return nil, errors.New("invalid status value")
		}
		task.Status = *req.Status
	}

	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

func (s *Service) DeleteTask(taskID, userID uuid.UUID) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to find task by ID: %w", err)
	}
	if task == nil {
		return domain.ErrTaskNotFound
	}

	if task.UserID != userID {
		return domain.ErrNoPermissionDeleteTask
	}

	if err := s.taskRepo.Delete(taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (s *Service) GetTaskByID(taskID, userID uuid.UUID) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task by ID: %w", err)
	}
	if task == nil {
		return nil, domain.ErrTaskNotFound
	}

	if task.UserID != userID {
		return nil, domain.ErrNoPermissionViewTask
	}

	return task, nil
}

//implement GetUserTasks function

// write out a proper status changing function and priority changing function

/*
func (s *Service) PendingTask(taskID, userID uuid.UUID) (*domain.Task, error) {
	pendingStatus := domain.TaskStatusPending
	return s.UpdateTask(taskID, userID, UpdateTaskRequest{
		Status: &pendingStatus,
	})
}

func (s *Service) InProgressTask(taskID, userID uuid.UUID) (*domain.Task, error) {
	inProgressStatus := domain.TaskStatusInProgress
	return s.UpdateTask(taskID, userID, UpdateTaskRequest{
		Status: &inProgressStatus,
	})
}

func (s *Service) CompletedTask(taskID, userID uuid.UUID) (*domain.Task, error) {
	completedStatus := domain.TaskStatusCompleted
	return s.UpdateTask(taskID, userID, UpdateTaskRequest{
		Status: &completedStatus,
	})
}

func (s *Service) CancelledTask(taskID, userID uuid.UUID) (*domain.Task, error) {
	cancelledStatus := domain.TaskStatusCancelled
	return s.UpdateTask(taskID, userID, UpdateTaskRequest{
		Status: &cancelledStatus,
	})
}
*/
