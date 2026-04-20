package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goobermv/calendar-task-tracker/internal/domain"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/dto"
	"github.com/goobermv/calendar-task-tracker/internal/network/api/middleware"
	taskUsecase "github.com/goobermv/calendar-task-tracker/internal/usescases/task"
)

type TaskHandler struct {
	taskService *taskUsecase.Service
}

func NewTaskHandler(taskService *taskUsecase.Service) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	useCaseReq := dto.CreateTaskRequest{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
	}

	task, err := h.taskService.CreateTask(useCaseReq)
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.TaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	c.JSON(http.StatusCreated, response)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	taskID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrTaskNotFound)
		return
	}

	task, err := h.taskService.GetTaskByID(taskID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	response := dto.TaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) GetUserTasks(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	tasks, err := h.taskService.GetUserTasks(userID)
	if err != nil {
		HandleError(c, domain.ErrFailedToGetUserTasks)
	}

	responses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = dto.TaskResponse{
			ID:          task.ID,
			UserID:      task.UserID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			DueDate:     task.DueDate,
			Priority:    task.Priority,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": responses,
		"count": len(responses),
	})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	taskID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrTaskNotFound)
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	task, err := h.taskService.UpdateTask(taskID, userID, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	resonse := dto.TaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	c.JSON(http.StatusOK, resonse)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		HandleError(c, domain.ErrUnauthorized)
		return
	}

	taskID, exists := middleware.GetID(c)
	if !exists {
		HandleError(c, domain.ErrTaskNotFound)
		return
	}

	err := h.taskService.DeleteTask(taskID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
